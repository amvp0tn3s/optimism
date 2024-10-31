package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	opnode "github.com/ethereum-optimism/optimism/op-node"
	"github.com/ethereum-optimism/optimism/op-node/flags"
	"github.com/ethereum-optimism/optimism/op-node/node"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	opflags "github.com/ethereum-optimism/optimism/op-service/flags"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
)

const EnvVarPrefix = "REPLAY"

var (
	ReplayTargetBlock = &cli.StringFlag{
		Name:    "replay-target",
		Usage:   "Hash of the block to be replayed",
		EnvVars: prefixEnvVars("TARGET"),
	}
	NumReplayedBlocks = &cli.IntFlag{
		Name:    "num-blocks",
		Usage:   "Hash of the block to be replayed",
		EnvVars: prefixEnvVars("NUM_BLOCKS"),
		Value:   100,
	}
	L2SyncAddr = &cli.StringFlag{
		Name:     "l2.sync",
		Usage:    "Address of the trusted L2 Engine JSON-RPC endpoints to use to retrieve payload information from",
		EnvVars:  prefixEnvVars("L2_SYNC_RPC"),
		Category: flags.RollupCategory,
	}
	BuilderL2EngineAddr = &cli.StringFlag{
		Name:     "l2.builder",
		Usage:    "Address of L2 builder's Engine JSON-RPC endpoints to use (engine and eth namespace required)",
		EnvVars:  prefixEnvVars("BUILDER_L2_ENGINE_RPC"),
		Category: flags.RollupCategory,
	}
	ValidatorL2EngineAddr = &cli.StringFlag{
		Name:     "l2.validator",
		Usage:    "Address of the validators's L2 Engine JSON-RPC endpoints to use (engine and eth namespace required)",
		EnvVars:  prefixEnvVars("VALIDATOR_L2_ENGINE_RPC"),
		Category: flags.RollupCategory,
	}
	L2EngineJWTSecret = &cli.StringFlag{
		Name:        "l2.jwt-secret",
		Usage:       "Path to JWT secret key. Keys are 32 bytes, hex encoded in a file. A new key will be generated if the file is empty.",
		EnvVars:     prefixEnvVars("L2_ENGINE_AUTH"),
		Value:       "",
		Destination: new(string),
		Category:    flags.RollupCategory,
	}
)

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

var requiredFlags = []cli.Flag{
	// both builder and validator have to use the same JWT secret for now
	L2EngineJWTSecret,
	BuilderL2EngineAddr,
	ValidatorL2EngineAddr,
	L2SyncAddr,
}

var optionalFlags = []cli.Flag{
	NumReplayedBlocks,
	ReplayTargetBlock,
}

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

func NewL2EndpointConfig(ctx *cli.Context, log log.Logger, addrCLIOpt *cli.StringFlag) (*node.L2EndpointConfig, error) {
	l2Addr := ctx.String(addrCLIOpt.Name)
	fileName := ctx.String(L2EngineJWTSecret.Name)
	var secret [32]byte
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return nil, fmt.Errorf("file-name of jwt secret is empty")
	}
	if data, err := os.ReadFile(fileName); err == nil {
		jwtSecret := common.FromHex(strings.TrimSpace(string(data)))
		if len(jwtSecret) != 32 {
			return nil, fmt.Errorf("invalid jwt secret in path %s, not 32 hex-formatted bytes", fileName)
		}
		copy(secret[:], jwtSecret)
	} else {
		log.Warn("Failed to read JWT secret from file, generating a new one now. Configure L2 geth with --authrpc.jwt-secret=" + fmt.Sprintf("%q", fileName))
		if _, err := io.ReadFull(rand.Reader, secret[:]); err != nil {
			return nil, fmt.Errorf("failed to generate jwt secret: %w", err)
		}
		if err := os.WriteFile(fileName, []byte(hexutil.Encode(secret[:])), 0o600); err != nil {
			return nil, err
		}
	}

	return &node.L2EndpointConfig{
		L2EngineAddr:      l2Addr,
		L2EngineJWTSecret: secret,
	}, nil
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
	l2ValidatorCfg, err := NewL2EndpointConfig(cliCtx, logger, ValidatorL2EngineAddr)
	if err != nil {
		return nil, logger, fmt.Errorf("failed to load l2 endpoints info: %w", err)
	}
	l2BuilderCfg, err := NewL2EndpointConfig(cliCtx, logger, BuilderL2EngineAddr)
	if err != nil {
		return nil, logger, fmt.Errorf("failed to load l2 endpoints info: %w", err)
	}

	l2LiveAddr := cliCtx.String(L2SyncAddr.Name)
	target := cliCtx.String(ReplayTargetBlock.Name)
	var targetBlock *common.Hash
	if target != "" {
		t := common.HexToHash(target)
		targetBlock = &t
	}
	numReplayBlocks := cliCtx.Int(NumReplayedBlocks.Name)
	return &Config{
		L2Validator:          l2ValidatorCfg,
		L2Builder:            l2BuilderCfg,
		L2Live:               l2LiveAddr,
		Rollup:               *rollupCfg,
		LogConfig:            logCfg,
		ReplayUntilBlockHash: targetBlock,
		NumReplayBlocks:      uint(numReplayBlocks),
	}, logger, nil
}
