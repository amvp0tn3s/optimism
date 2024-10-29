package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	L2        node.L2EndpointSetup
	Rollup    rollup.Config
	LogConfig oplog.CLIConfig

	// TODO: have as arg to a "replay" cmd
	ReplayBlockHash common.Hash
}

func (cfg *Config) Check() error {
	if err := cfg.Rollup.Check(); err != nil {
		return err
	}
	if err := cfg.L2.Check(); err != nil {
		return err
	}
	// TODO: check blockhash
	return nil
}

type ConsensusOrchestrator struct {
	running atomic.Bool
	client  *sources.EngineClient
	config  *Config
	logger  log.Logger

	errs <-chan error
}

func NewConsensusOrchestrator(ctx context.Context, config *Config, logger log.Logger) (*ConsensusOrchestrator, error) {
	var log log.Logger

	rpcClient, rpcCfg, err := config.L2.Setup(ctx, log, &config.Rollup)
	if err != nil {
		return nil, fmt.Errorf("failed to setup L2 execution-engine RPC client: %w", err)
	}
	engineClient, err := sources.NewEngineClient(
		rpcClient, log, nil, rpcCfg,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate engine client: %w", err)
	}
	return &ConsensusOrchestrator{
		client: engineClient,
		config: config,
		logger: logger,
	}, nil
}

func (co *ConsensusOrchestrator) Start(ctx context.Context) error {
	co.logger.Info("Starting consensus orchestrator")
	co.running.Store(true)
	c := make(chan error, 10)
	// TODO: use errorgroup
	co.errs = c
	go co.replay(ctx, c, co.config.ReplayBlockHash)
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
