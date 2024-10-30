package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	opnode "github.com/ethereum-optimism/optimism/op-node"
	"github.com/ethereum-optimism/optimism/op-node/flags"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	opflags "github.com/ethereum-optimism/optimism/op-service/flags"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
)

const EnvVarPrefix = "REPLAY"

var (
	ReplayBlock = &cli.StringFlag{
		Name:    "replay-block",
		Usage:   "Hash of the block to be replayed",
		EnvVars: prefixEnvVars("BLOCK"),
	}
	L2SyncAddr = &cli.StringFlag{
		Name:     "l2-sync",
		Usage:    "Address of the trusted L2 Engine JSON-RPC endpoints to use to retrieve payload information from",
		EnvVars:  prefixEnvVars("L2_SYNC_RPC"),
		Category: flags.RollupCategory,
	}
)

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

var requiredFlags = []cli.Flag{
	flags.L2EngineAddr,
	flags.L2EngineJWTSecret,
	L2SyncAddr,

	// TODO: I'd rather have this as an argument
	ReplayBlock,
}

var optionalFlags = []cli.Flag{}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func init() {
	requiredFlags = append(requiredFlags, opflags.CLIRollupConfigFlag(flags.EnvVarPrefix, flags.RollupCategory))
	optionalFlags = append(optionalFlags, oplog.CLIFlags(EnvVarPrefix)...)
	Flags = append(requiredFlags, optionalFlags...)
}

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}

func ParseConfig(cliCtx *cli.Context) (*Config, log.Logger, error) {
	logCfg := oplog.ReadCLIConfig(cliCtx)
	logger := oplog.NewLogger(oplog.AppOut(cliCtx), logCfg)

	err := CheckRequired(cliCtx)
	if err != nil {
		return nil, logger, err
	}
	rollupCfg, err := opnode.NewRollupConfigFromCLI(logger, cliCtx)
	if err != nil {
		return nil, logger, fmt.Errorf("failed to build rollup config: %w", err)
	}
	// HACK: overwrite this so that it forces an EngineAPI V2...
	rollupCfg.GraniteTime = nil
	rollupCfg.FjordTime = nil
	rollupCfg.EcotoneTime = nil
	l2Cfg, err := opnode.NewL2EndpointConfig(cliCtx, logger)
	if err != nil {
		return nil, logger, fmt.Errorf("failed to load l2 endpoints info: %w", err)
	}

	l2LiveAddr := cliCtx.String(L2SyncAddr.Name)
	block := common.HexToHash(cliCtx.String(ReplayBlock.Name))
	return &Config{
		L2Local:         l2Cfg,
		L2Live:          l2LiveAddr,
		Rollup:          *rollupCfg,
		LogConfig:       logCfg,
		ReplayBlockHash: block,
	}, logger, nil
}
