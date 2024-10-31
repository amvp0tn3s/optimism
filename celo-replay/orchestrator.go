package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	L2Validator node.L2EndpointSetup
	L2Builder   node.L2EndpointSetup
	L2Live      string
	Rollup      rollup.Config
	LogConfig   oplog.CLIConfig

	ReplayUntilBlockHash *common.Hash
	NumReplayBlocks      uint
}

func (cfg *Config) Check() error {
	if err := cfg.Rollup.Check(); err != nil {
		return err
	}
	if err := cfg.L2Validator.Check(); err != nil {
		return err
	}
	// TODO: check blockhash
	return nil
}

type SyncStatus struct {
	Status eth.PayloadStatusV1
	Error  error

	Replay          bool
	Building        bool
	BuildAttributes *eth.PayloadAttributes
	BuildingID      *eth.PayloadID
	BuiltPayload    *eth.ExecutionPayloadEnvelope
	TargetPayload   *eth.ExecutionPayloadEnvelope
}

func (ss *SyncStatus) SyncedHash() *common.Hash {
	return ss.Status.LatestValidHash
}

// func (ss *SyncStatus) SyncedNumber() *eth.Uint64Quantity {
// 	return nil
// }

type ConsensusOrchestrator struct {
	running        atomic.Bool
	validateClient *EngineClient
	buildClient    *EngineClient
	liveClient     *sources.L2Client
	config         *Config
	logger         log.Logger
}

type EngineClient struct {
	*sources.EngineClient
	Name string
}

func SetupEngineClient(
	ctx context.Context,
	endpointSetup node.L2EndpointSetup,
	name string,
	config *Config,
	logger log.Logger,
) (*EngineClient, error) {
	rpcClient, rpcCfg, err := endpointSetup.Setup(ctx, logger, &config.Rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to setup L2 %s execution-engine RPC client: %w", name, err)
	}
	builderEngineClient, err := sources.NewEngineClient(
		rpcClient, logger, nil, rpcCfg,
	)
	return &EngineClient{
		EngineClient: builderEngineClient,
		Name:         name,
	}, nil
}
func NewConsensusOrchestrator(ctx context.Context, config *Config, logger log.Logger) (*ConsensusOrchestrator, error) {
	validatorEngineClient, err := SetupEngineClient(ctx, config.L2Validator, "validator", config, logger)
	if err != nil {
		return nil, err
	}
	builderEngineClient, err := SetupEngineClient(ctx, config.L2Builder, "builder", config, logger)
	if err != nil {
		return nil, err
	}

	liveClient, err := setupL2LiveClient(ctx, config, logger)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate live client: %w", err)
	}
	return &ConsensusOrchestrator{
		validateClient: validatorEngineClient,
		buildClient:    builderEngineClient,
		liveClient:     liveClient,
		config:         config,
		logger:         logger,
	}, nil
}

func setupL2LiveClient(ctx context.Context, config *Config, logger log.Logger) (*sources.L2Client, error) {
	opts := []client.RPCOption{
		client.WithDialBackoff(10),
	}
	l2Node, err := client.NewRPC(ctx, logger, config.L2Live, opts...)
	if err != nil {
		return nil, err
	}
	liveClient, err := sources.NewL2Client(
		l2Node, logger, nil, sources.L2ClientDefaultConfig(&config.Rollup, true),
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate live syncing client: %w", err)
	}
	return liveClient, nil
}

func (co *ConsensusOrchestrator) Start(ctx context.Context) error {
	co.logger.Info("Starting consensus orchestrator")
	co.running.Store(true)
	co.run(ctx)
	return nil
}
func (co *ConsensusOrchestrator) Stop(ctx context.Context) error {
	if !co.running.Load() {
		return errors.New("already stopped")
	}
	co.logger.Info("Stopping consensus orchestrator")
	co.running.Store(false)
	return nil
}
func (co *ConsensusOrchestrator) Stopped() bool {
	return co.running.Load()
}

func (co *ConsensusOrchestrator) monitor(ctx context.Context, s <-chan SyncStatus) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case status, ok := <-s:
			if !ok {
				// channel closed, there is nothing else to monitor
				return nil
			}
			if status.Replay {
				if status.Error != nil {
					co.logger.Error("error during rebuilding", "error", status.Error)
				}
				if status.Status.ValidationError != nil {
					co.logger.Error("validation error during block rebuild",
						"error", status.Status.ValidationError,
						//FIXME: could be nil deref
						"built-hash", status.BuiltPayload.ExecutionPayload.BlockHash,
						"built-number", uint64(status.BuiltPayload.ExecutionPayload.BlockNumber),

						"target-hash", status.TargetPayload.ExecutionPayload.BlockHash,
						"target-number", uint64(status.TargetPayload.ExecutionPayload.BlockNumber),
						"latest-valid-hash", status.Status.LatestValidHash,
						"status", status.Status.Status,
					)
					// TODO: add to a hashmap,
					// so we can check if they later got validated during FCU
					// with the original payload,
					// and wether there is a payload difference or
					// not (race-condition!)
				} else {
					co.logger.Info("re-built old block",
						"built-hash", status.BuiltPayload.ExecutionPayload.BlockHash,
						"built-number", uint64(status.BuiltPayload.ExecutionPayload.BlockNumber),
						"target-hash", status.TargetPayload.ExecutionPayload.BlockHash,
						"target-number", uint64(status.TargetPayload.ExecutionPayload.BlockNumber),
						"latest-valid-hash", status.Status.LatestValidHash,
						"status", status.Status.Status,
					)
				}
			} else {
				if status.Error != nil {
					co.logger.Error("sync status error", "error", status.Error, "block", status.TargetPayload.ExecutionPayload.BlockHash)
				}
				co.logger.Info("synced fork-choice",
					"hash", status.TargetPayload.ExecutionPayload.BlockHash,
					"number", uint64(status.TargetPayload.ExecutionPayload.BlockNumber),
					"latest-valid-hash", status.Status.LatestValidHash,
					"status", status.Status.Status,
				)
			}
		}
	}
}
