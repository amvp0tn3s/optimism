package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"runtime/debug"
	"time"

	"log/slog"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-service/ioutil"
	"github.com/ethereum-optimism/optimism/op-service/jsonutil"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/mattn/go-isatty"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"golang.org/x/sync/errgroup"
)

var (
	deployConfigFlag = &cli.PathFlag{
		Name:     "deploy-config",
		Usage:    "Path to the JSON file that was used for the bedrock contracts deployment. A test example can be found here 'op-chain-ops/genesis/testdata/test-deploy-config-full.json' and documentation for the fields is at https://docs.optimism.io/builders/chain-operators/management/configuration",
		Required: true,
	}
	l1DeploymentsFlag = &cli.PathFlag{
		Name:     "l1-deployments",
		Usage:    "Path to L1 deployments JSON file, the output of running the bedrock contracts deployment for the given 'deploy-config'",
		Required: true,
	}
	l1RPCFlag = &cli.StringFlag{
		Name:     "l1-rpc",
		Usage:    "RPC URL for a node of the L1 defined in the 'deploy-config'",
		Required: true,
	}
	l2AllocsFlag = &cli.PathFlag{
		Name:     "l2-allocs",
		Usage:    "Path to L2 genesis allocs file. You can find instructions on how to generate this file in the README",
		Required: true,
	}
	outfileRollupConfigFlag = &cli.PathFlag{
		Name:     "outfile.rollup-config",
		Usage:    "Path to write the rollup config JSON file, to be provided to op-node with the 'rollup.config' flag",
		Required: true,
	}
	outfileGenesisFlag = &cli.PathFlag{
		Name:     "outfile.genesis",
		Usage:    "Path to write the genesis JSON file, to be used to sync new nodes",
		Required: true,
	}
	migrationBlockNumberFlag = &cli.Uint64Flag{
		Name:     "migration-block-number",
		Usage:    "Specifies the migration block number. If the source db is not synced exactly to the block immediately before this number (i.e. migration-block-number - 1), the migration will fail.",
		Required: true,
	}
	migrationBlockTimeFlag = &cli.Uint64Flag{
		Name:     "migration-block-time",
		Usage:    "Specifies a unix timestamp to use for the migration block. This should be set to the same timestamp as was used for the sequencer migration. If performing the sequencer migration, this should be set to a time in the future around when the migration script is expected to complete.",
		Required: true,
	}
	oldDBPathFlag = &cli.PathFlag{
		Name:     "old-db",
		Usage:    "Path to the old Celo chaindata dir, can be found at '<datadir>/celo/chaindata'",
		Required: true,
	}
	newDBPathFlag = &cli.PathFlag{
		Name:     "new-db",
		Usage:    "Path to write migrated Celo chaindata, note the new node implementation expects to find this chaindata at the following path '<datadir>/geth/chaindata",
		Required: true,
	}
	batchSizeFlag = &cli.Uint64Flag{
		Name:  "batch-size",
		Usage: "Batch size to use for block migration, larger batch sizes can speed up migration but require more memory. If increasing the batch size consider also increasing the memory-limit",
		Value: 50000, // TODO(Alec) optimize default parameters
	}
	bufferSizeFlag = &cli.Uint64Flag{
		Name:  "buffer-size",
		Usage: "Buffer size to use for ancient block migration channels. Defaults to 0. Included to facilitate testing for performance improvements.",
		Value: 0,
	}
	memoryLimitFlag = &cli.Int64Flag{
		Name:  "memory-limit",
		Usage: "Memory limit in MiB, should be set lower than the available amount of memory in your system to prevent out of memory errors",
		Value: 7500,
	}
	reset = &cli.BoolFlag{
		Name:  "reset",
		Usage: "Delete everything in the destination directory aside from /ancients. This is useful if you need to re-run the full migration but do not want to repeat the lengthy ancients migration. If you'd like to reset the entire destination directory, you can delete it manually.",
		Value: false,
	}

	preMigrationFlags = []cli.Flag{
		oldDBPathFlag,
		newDBPathFlag,
		batchSizeFlag,
		bufferSizeFlag,
		memoryLimitFlag,
		reset,
	}
	fullMigrationFlags = append(
		preMigrationFlags,
		deployConfigFlag,
		l1DeploymentsFlag,
		l1RPCFlag,
		l2AllocsFlag,
		outfileRollupConfigFlag,
		outfileGenesisFlag,
		migrationBlockTimeFlag,
		migrationBlockNumberFlag,
	)
)

type preMigrationOptions struct {
	oldDBPath        string
	newDBPath        string
	batchSize        uint64
	bufferSize       uint64
	memoryLimit      int64
	resetNonAncients bool
}

type stateMigrationOptions struct {
	deployConfig        string
	l1Deployments       string
	l1RPC               string
	l2AllocsPath        string
	outfileRollupConfig string
	outfileGenesis      string
	migrationBlockTime  uint64
}

type fullMigrationOptions struct {
	preMigrationOptions
	stateMigrationOptions
	migrationBlockNumber uint64
}

func parsePreMigrationOptions(ctx *cli.Context) preMigrationOptions {
	return preMigrationOptions{
		oldDBPath:        ctx.String(oldDBPathFlag.Name),
		newDBPath:        ctx.String(newDBPathFlag.Name),
		batchSize:        ctx.Uint64(batchSizeFlag.Name),
		bufferSize:       ctx.Uint64(bufferSizeFlag.Name),
		memoryLimit:      ctx.Int64(memoryLimitFlag.Name),
		resetNonAncients: ctx.Bool(reset.Name),
	}
}

func parseStateMigrationOptions(ctx *cli.Context) stateMigrationOptions {
	return stateMigrationOptions{
		deployConfig:        ctx.Path(deployConfigFlag.Name),
		l1Deployments:       ctx.Path(l1DeploymentsFlag.Name),
		l1RPC:               ctx.String(l1RPCFlag.Name),
		l2AllocsPath:        ctx.Path(l2AllocsFlag.Name),
		outfileRollupConfig: ctx.Path(outfileRollupConfigFlag.Name),
		outfileGenesis:      ctx.Path(outfileGenesisFlag.Name),
		migrationBlockTime:  ctx.Uint64(migrationBlockTimeFlag.Name),
	}
}

func parseFullMigrationOptions(ctx *cli.Context) fullMigrationOptions {
	return fullMigrationOptions{
		preMigrationOptions:   parsePreMigrationOptions(ctx),
		stateMigrationOptions: parseStateMigrationOptions(ctx),
		migrationBlockNumber:  ctx.Uint64(migrationBlockNumberFlag.Name),
	}
}

func main() {

	color := isatty.IsTerminal(os.Stderr.Fd())
	handler := log.NewTerminalHandlerWithLevel(os.Stderr, slog.LevelInfo, color)
	oplog.SetGlobalLogHandler(handler)

	app := &cli.App{
		Name:  "celo-migrate",
		Usage: "Migrate Celo block and state data to a CeL2 DB",
		Commands: []*cli.Command{
			{
				Name:  "pre",
				Usage: "Perform a pre-migration of ancient blocks and copy over all other data without transforming it. This should be run a day before the full migration command is run to minimize downtime.",
				Flags: preMigrationFlags,
				Action: func(ctx *cli.Context) error {
					if _, _, err := runPreMigration(parsePreMigrationOptions(ctx)); err != nil {
						return fmt.Errorf("failed to run pre-migration: %w", err)
					}
					log.Info("Finished pre-migration successfully!")
					return nil
				},
			},
			{
				Name:  "full",
				Usage: "Perform a full migration of both block and state data to a CeL2 DB",
				Flags: fullMigrationFlags,
				Action: func(ctx *cli.Context) error {
					if err := runFullMigration(parseFullMigrationOptions(ctx)); err != nil {
						return fmt.Errorf("failed to run full migration: %w", err)
					}
					log.Info("Finished full migration successfully!")
					return nil
				},
			},
		},
		OnUsageError: func(ctx *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}
			_ = cli.ShowAppHelp(ctx)
			return fmt.Errorf("please provide a valid command")
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Crit("error in migration", "err", err)
	}
}

func runFullMigration(opts fullMigrationOptions) error {
	defer timer("full migration")()

	log.Info("Full Migration Started", "oldDBPath", opts.oldDBPath, "newDBPath", opts.newDBPath)

	head, err := getHeadHeader(opts.oldDBPath)
	if err != nil {
		return fmt.Errorf("failed to get head header: %w", err)
	}
	if head.Number.Uint64() != opts.migrationBlockNumber-1 {
		return fmt.Errorf("old-db head block number not synced to the block immediately before the migration block number: %d != %d", head.Number.Uint64(), opts.migrationBlockNumber-1)
	}

	log.Info("Source db is synced to correct height", "head", head.Number.Uint64(), "migrationBlock", opts.migrationBlockNumber)

	var numAncients uint64
	var strayAncientBlocks []*rawdb.NumberHash

	if strayAncientBlocks, numAncients, err = runPreMigration(opts.preMigrationOptions); err != nil {
		return fmt.Errorf("failed to run pre-migration: %w", err)
	}

	if err = runNonAncientMigration(opts.newDBPath, strayAncientBlocks, opts.batchSize, numAncients); err != nil {
		return fmt.Errorf("failed to run non-ancient migration: %w", err)
	}
	if err = runStateMigration(opts.newDBPath, opts.stateMigrationOptions); err != nil {
		return fmt.Errorf("failed to run state migration: %w", err)
	}

	log.Info("Full Migration Finished", "oldDBPath", opts.oldDBPath, "newDBPath", opts.newDBPath)

	return nil
}

func runPreMigration(opts preMigrationOptions) ([]*rawdb.NumberHash, uint64, error) {
	defer timer("pre-migration")()

	log.Info("Pre-Migration Started", "oldDBPath", opts.oldDBPath, "newDBPath", opts.newDBPath, "batchSize", opts.batchSize, "memoryLimit", opts.memoryLimit)

	// Check that `rsync` command is available. We use this to copy the db excluding ancients, which we will copy separately
	if _, err := exec.LookPath("rsync"); err != nil {
		return nil, 0, fmt.Errorf("please install `rsync` to run block migration")
	}

	debug.SetMemoryLimit(opts.memoryLimit * 1 << 20) // Set memory limit, converting from MiB to bytes

	var err error

	if err = createNewDbPathIfNotExists(opts.newDBPath); err != nil {
		return nil, 0, fmt.Errorf("failed to create new db path: %w", err)
	}

	if opts.resetNonAncients {
		if err = cleanupNonAncientDb(opts.newDBPath); err != nil {
			return nil, 0, fmt.Errorf("failed to cleanup non-ancient db: %w", err)
		}
	}

	var numAncientsNewBefore uint64
	var numAncientsNewAfter uint64
	var strayAncientBlocks []*rawdb.NumberHash
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		if numAncientsNewBefore, numAncientsNewAfter, err = migrateAncientsDb(ctx, opts.oldDBPath, opts.newDBPath, opts.batchSize, opts.bufferSize); err != nil {
			return fmt.Errorf("failed to migrate ancients database: %w", err)
		}
		// Scanning for stray ancient blocks is slow, so we do it as soon as we can after the lock on oldDB is released by migrateAncientsDb
		// Doing this in parallel with copyDbExceptAncients still saves time if ancients have already been pre-migrated
		if strayAncientBlocks, err = getStrayAncientBlocks(opts.oldDBPath); err != nil {
			return fmt.Errorf("failed to get stray ancient blocks: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		// By doing this once during the premigration, we get a speedup when we run it again in a full migration.
		return copyDbExceptAncients(opts.oldDBPath, opts.newDBPath)
	})

	if err = g.Wait(); err != nil {
		return nil, 0, fmt.Errorf("failed to migrate blocks: %w", err)
	}

	log.Info("Pre-Migration Finished", "oldDBPath", opts.oldDBPath, "newDBPath", opts.newDBPath, "migratedAncients", numAncientsNewAfter-numAncientsNewBefore, "strayAncientBlocks", len(strayAncientBlocks))

	return strayAncientBlocks, numAncientsNewAfter, nil
}

func runNonAncientMigration(newDBPath string, strayAncientBlocks []*rawdb.NumberHash, batchSize, numAncients uint64) error {
	defer timer("non-ancient migration")()

	newDB, err := openDBWithoutFreezer(newDBPath, false)
	if err != nil {
		return fmt.Errorf("failed to open new database: %w", err)
	}
	defer newDB.Close()

	// get the last block number
	hash := rawdb.ReadHeadHeaderHash(newDB)
	lastBlock := *rawdb.ReadHeaderNumber(newDB, hash)
	lastAncient := numAncients - 1

	log.Info("Non-Ancient Block Migration Started", "process", "non-ancients", "newDBPath", newDBPath, "batchSize", batchSize, "startBlock", numAncients, "endBlock", lastBlock, "count", lastBlock-lastAncient, "lastAncientBlock", lastAncient)

	var numNonAncients uint64
	if numNonAncients, err = migrateNonAncientsDb(newDB, lastBlock, numAncients, batchSize); err != nil {
		return fmt.Errorf("failed to migrate non-ancients database: %w", err)
	}

	err = removeBlocks(newDB, strayAncientBlocks)
	if err != nil {
		return fmt.Errorf("failed to remove stray ancient blocks: %w", err)
	}
	log.Info("Removed stray ancient blocks still in leveldb", "process", "non-ancients", "removedBlocks", len(strayAncientBlocks))

	log.Info("Non-Ancient Block Migration Completed", "process", "non-ancients", "migratedNonAncients", numNonAncients)

	return nil
}

func runStateMigration(newDBPath string, opts stateMigrationOptions) error {
	defer timer("state migration")()

	log.Info("State Migration Started", "newDBPath", newDBPath, "deployConfig", opts.deployConfig, "l1Deployments", opts.l1Deployments, "l1RPC", opts.l1RPC, "l2AllocsPath", opts.l2AllocsPath, "outfileRollupConfig", opts.outfileRollupConfig)

	// Read deployment configuration
	config, err := genesis.NewDeployConfig(opts.deployConfig)
	if err != nil {
		return err
	}

	if config.DeployCeloContracts {
		return errors.New("DeployCeloContracts is not supported in migration")
	}
	if config.FundDevAccounts {
		return errors.New("FundDevAccounts is not supported in migration")
	}

	// Try reading the L1 deployment information
	deployments, err := genesis.NewL1Deployments(opts.l1Deployments)
	if err != nil {
		return fmt.Errorf("cannot read L1 deployments at %s: %w", opts.l1Deployments, err)
	}
	config.SetDeployments(deployments)

	// Get latest block information from L1
	var l1StartBlock *types.Block
	client, err := ethclient.Dial(opts.l1RPC)
	if err != nil {
		return fmt.Errorf("cannot dial %s: %w", opts.l1RPC, err)
	}

	if config.L1StartingBlockTag == nil {
		l1StartBlock, err = client.BlockByNumber(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("cannot fetch latest block: %w", err)
		}
		tag := rpc.BlockNumberOrHashWithHash(l1StartBlock.Hash(), true)
		config.L1StartingBlockTag = (*genesis.MarshalableRPCBlockNumberOrHash)(&tag)
	} else if config.L1StartingBlockTag.BlockHash != nil {
		l1StartBlock, err = client.BlockByHash(context.Background(), *config.L1StartingBlockTag.BlockHash)
		if err != nil {
			return fmt.Errorf("cannot fetch block by hash: %w", err)
		}
	} else if config.L1StartingBlockTag.BlockNumber != nil {
		l1StartBlock, err = client.BlockByNumber(context.Background(), big.NewInt(config.L1StartingBlockTag.BlockNumber.Int64()))
		if err != nil {
			return fmt.Errorf("cannot fetch block by number: %w", err)
		}
	}

	// Ensure that there is a starting L1 block
	if l1StartBlock == nil {
		return fmt.Errorf("no starting L1 block")
	}

	// Sanity check the config. Do this after filling in the L1StartingBlockTag
	// if it is not defined.
	if err := config.Check(log.New()); err != nil {
		return err
	}

	log.Info("Using L1 Start Block", "number", l1StartBlock.Number(), "hash", l1StartBlock.Hash().Hex())

	// Build the L2 genesis block
	l2Allocs, err := foundry.LoadForgeAllocs(opts.l2AllocsPath)
	if err != nil {
		return err
	}

	l2Genesis, err := genesis.BuildL2Genesis(config, l2Allocs, l1StartBlock.Header())
	if err != nil {
		return fmt.Errorf("error creating l2 genesis: %w", err)
	}

	// Write changes to state to actual state database
	cel2Header, err := applyStateMigrationChanges(config, l2Genesis.Alloc, newDBPath, opts.outfileGenesis, opts.migrationBlockTime, l1StartBlock)
	if err != nil {
		return err
	}
	log.Info("Updated Cel2 state")

	rollupConfig, err := config.RollupConfig(l1StartBlock.Header(), cel2Header.Hash(), cel2Header.Number.Uint64())
	if err != nil {
		return err
	}
	if err := rollupConfig.Check(); err != nil {
		return fmt.Errorf("generated rollup config does not pass validation: %w", err)
	}

	log.Info("Writing rollup config", "file", opts.outfileRollupConfig)
	if err := jsonutil.WriteJSON(rollupConfig, ioutil.ToStdOutOrFileOrNoop(opts.outfileRollupConfig, OutFilePerm)); err != nil {
		return err
	}

	log.Info("State Migration Completed")

	return nil
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Info("TIMER", "process", name, "duration", time.Since(start))
	}
}
