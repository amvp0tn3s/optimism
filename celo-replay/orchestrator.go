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
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	L2Local   node.L2EndpointSetup
	L2Live    string
	Rollup    rollup.Config
	LogConfig oplog.CLIConfig

	// TODO: have as arg to a "replay" cmd
	ReplayBlockHash common.Hash
}

func (cfg *Config) Check() error {
	if err := cfg.Rollup.Check(); err != nil {
		return err
	}
	if err := cfg.L2Local.Check(); err != nil {
		return err
	}
	// TODO: check blockhash
	return nil
}

type ConsensusOrchestrator struct {
	running    atomic.Bool
	client     *sources.EngineClient
	liveClient *sources.L2Client
	config     *Config
	logger     log.Logger

	errs <-chan error
}

func NewConsensusOrchestrator(ctx context.Context, config *Config, logger log.Logger) (*ConsensusOrchestrator, error) {
	// Local engine client setup
	rpcClient, rpcCfg, err := config.L2Local.Setup(ctx, logger, &config.Rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to setup L2 execution-engine RPC client: %w", err)
	}

	// Use a static version 2 for now
	// v2Provider := engine.StaticVersionProvider(2)
	// rpcCfg.RollupCfg = v2Provider
	engineClient, err := sources.NewEngineClient(
		rpcClient, logger, nil, rpcCfg,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate engine client: %w", err)
	}

	liveClient, err := setupL2LiveClient(ctx, config, logger)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate live client: %w", err)
	}
	return &ConsensusOrchestrator{
		client:     engineClient,
		liveClient: liveClient,
		config:     config,
		logger:     logger,
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
	c := make(chan error, 10)
	// TODO: use errorgroup
	co.errs = c
	go co.forkChoice(ctx, c, co.config.ReplayBlockHash)
	go co.monitor(ctx)
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

func (co *ConsensusOrchestrator) monitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-co.errs:
			co.logger.Crit("", "error", err)
			return
		}
	}
}
