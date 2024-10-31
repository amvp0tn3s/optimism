package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

var (
	errExecutionInvalid = errors.New("execution invalid")
)

func (co *ConsensusOrchestrator) newPayload(ctx context.Context, engine *EngineClient, payload *eth.ExecutionPayloadEnvelope) (*eth.PayloadStatusV1, error) {
	status, err := engine.EngineAPIClient.NewPayload(ctx, payload.ExecutionPayload, payload.ParentBeaconBlockRoot)
	if err != nil {
		return status, fmt.Errorf("failed to send newPaylod to %s client: %w", engine.Name, err)
	}
	if status.Status == eth.ExecutionInvalid {
		err := errors.New(*status.ValidationError)
		return status, fmt.Errorf("error on %s client: %w: %w", engine.Name, errExecutionInvalid, err)
	}
	return status, nil
}

func (co *ConsensusOrchestrator) syncUntil(
	ctx context.Context,
	s chan<- SyncStatus,
	engine *EngineClient,
	target *eth.ExecutionPayloadEnvelope,
	buildAttrs *eth.PayloadAttributes,
) (*SyncStatus, error) {
	if target == nil {
		return nil, errors.New("nil target not allowed")
	}
	// first get the checkpoint payload and forward it to the EL
	co.logger.Debug("setting target payload", "block", target.ExecutionPayload.BlockHash.Hex())
	status, err := co.newPayload(ctx, engine, target)
	if err != nil {
		sstatus := &SyncStatus{
			Status:        *status,
			Error:         err,
			TargetPayload: target,
		}
		// got an invalid block
		s <- *sstatus
		return sstatus, err
	}
	// the EL has the payload now, we can start syncing:
	t := time.NewTicker(5 * time.Second)
	trigger := make(chan struct{}, 1)
	setTrigger := func() {
		select {
		case trigger <- struct{}{}:
		// if this is block we are triggering
		// anyways
		default:
			return
		}
	}
	defer close(trigger)
	defer t.Stop()

	setTrigger()
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-t.C:
			setTrigger()
		case <-trigger:
			status, payloadId, err := co.forkChoice(ctx, target.ExecutionPayload.BlockHash, buildAttrs)
			sstatus := &SyncStatus{
				Status:          *status,
				Building:        bool(payloadId != nil),
				BuildingID:      payloadId,
				BuildAttributes: buildAttrs,
				Error:           err,
				TargetPayload:   target,
			}
			co.logger.Trace("called FCU", "status", sstatus)
			s <- *sstatus
			if err != nil {
				return sstatus, err
			}
			if status.Status == eth.ExecutionSyncing {
				co.logger.Debug("execution syncing after FCU", "hash", target.ExecutionPayload.BlockHash)
				// do nothing and continue to call the same fork-choice until we are synced
				// or hit e.g. an invalid block
				continue
			}
			if status.Status == eth.ExecutionValid {
				co.logger.Debug("reached valid target FCU", "hash", target.ExecutionPayload.BlockHash)
				return sstatus, err
			}
		}

	}

}

func (co *ConsensusOrchestrator) run(ctx context.Context) error {
	st := make(chan SyncStatus, 10)
	errg, ctx := errgroup.WithContext(ctx)

	errg.Go(func() error {
		defer close(st)
		return co.replayer(ctx, st)
	})
	errg.Go(func() error {
		return co.monitor(ctx, st)
	})
	go co.waitShutdown(errg)
	return nil
}
func (co *ConsensusOrchestrator) waitShutdown(errgrp *errgroup.Group) {
	err := errgrp.Wait()
	if err != nil {
		co.logger.Crit("received error, shutting down", "error", err)
	}
	os.Exit(0)
}

func (co *ConsensusOrchestrator) getPayloadChain(ctx context.Context, latest common.Hash, num uint) ([]*eth.ExecutionPayloadEnvelope, error) {
	payloads := make([]*eth.ExecutionPayloadEnvelope, 0)
	currentHash := latest
	initial := num
	n := num
	for n != 0 {
		if n == 1 || bool((n/initial)%10 == 0) {
			co.logger.Info("fetching block payloads from trusted RPC server",
				"remaining", n-initial,
				"total", initial,
			)
		}
		payload, err := co.liveClient.PayloadByHash(ctx, currentHash)
		if err != nil {
			return nil, fmt.Errorf("couldn't retrieve payload from trusted RPC provider: %w", err)
		}
		//FIXME: payload nil check?
		if payload == nil {
			return nil, errors.New("nil payload received")
		}
		payloads = append(payloads, payload)
		currentHash = payload.ExecutionPayload.ParentHash
		n--

		// don't overload the rpc
		// OPTIM: arbitrary wait time set, investigate what would work better
		time.Sleep(300 * time.Millisecond)
	}
	return payloads, nil
}
func (co *ConsensusOrchestrator) replayer(ctx context.Context, s chan<- SyncStatus) error {

	if co.config.ReplayUntilBlockHash == nil {
		target, err := co.liveClient.PayloadByLabel(ctx, eth.Unsafe)
		if err != nil {
			return fmt.Errorf("can't fetch latest unsafe head from trusted RPC server: %w ", err)
		}
		co.config.ReplayUntilBlockHash = &target.ExecutionPayload.BlockHash
		co.logger.Info("no target replay hash was provided, using latest unsafe",
			"number", uint64(target.ExecutionPayload.BlockNumber),
			"hash", target.ExecutionPayload.BlockHash,
		)
	}
	co.logger.Info("pre-fetching all block payloads from trusted RPC server", "num-blocks", co.config.NumReplayBlocks)
	if co.config.NumReplayBlocks > 10000 {
		co.logger.Crit("number of prefetch blocks is too high, the orchestrator is not optimised for this yet",
			"num-blocks", co.config.NumReplayBlocks,
		)
		return errors.New("number fetch blocks is too high")
	}
	if co.config.NumReplayBlocks > 1000 {
		co.logger.Warn("number of prefetch blocks is high, this will take some time", "num-blocks", co.config.NumReplayBlocks)
	}
	// OPTIM: if long block-ranges should be re-synced, this will take some time
	// and load everything into memory in one go. We could do this in an ad-hoc fashion
	// and only download chunks
	payloads, err := co.getPayloadChain(ctx, *co.config.ReplayUntilBlockHash, co.config.NumReplayBlocks+2)
	if err != nil {
		return err
	}
	finalPayload := payloads[0]
	// we need this only for the last timestamp
	payloadAfterFinal, err := co.liveClient.PayloadByNumber(ctx, uint64(finalPayload.ExecutionPayload.BlockNumber))
	if err != nil {
		return fmt.Errorf("couldn't retrieve payload from trusted RPC provider: %w", err)
	}
	beforeBuildingPayload := payloads[len(payloads)-1]
	payloadIndex := uint(len(payloads) - 1)
	building := false
	var buildAttrs *eth.PayloadAttributes
	co.logger.Info("syncing to before the replay section",
		"hash", beforeBuildingPayload.ExecutionPayload.BlockHash,
		"number", uint64(beforeBuildingPayload.ExecutionPayload.BlockNumber),
	)

	// TODO:
	// - query the current block-number and hash
	// - if we are after the beforeBuildingPayload block,
	//	 use something like the `debug_setHead` to rewind the chain-state instead first
	for {
		currentPayload := payloads[payloadIndex]
		if building && payloadIndex == 0 {
			building = false
			buildAttrs = nil
		}
		if building {
			buildPayload := payloads[payloadIndex-1]
			buildAttrs = PayloadToBuildAttrs(*buildPayload)
		}
		status, err := co.syncUntil(ctx, s, currentPayload, buildAttrs)
		if err != nil {
			co.logger.Error("syncer errored", "error", err)
			// TODO: we might try to recover from this
			return err
		}
		co.logger.Debug("successfully synced to current sync target", "status", status)
		if status.SyncedHash() == nil {
			// XXX: are there scenarios where this is expected?
			co.logger.Error("no synced hash returned from engine api", "error", err)
			return err
		}
		if status.SyncedHash().Cmp(finalPayload.ExecutionPayload.BlockHash) == 0 {
			co.logger.Info("successfully synced to final sync target", "block-hash", status.SyncedHash())
			break
		} else if status.SyncedHash().Cmp(beforeBuildingPayload.ExecutionPayload.BlockHash) == 0 {
			co.logger.Info("successfully synced to before building sync target", "block-hash", status.SyncedHash())
			// start building
			co.logger.Info("starting re-build and validating blocks until target",
				"num-blocks", co.config.NumReplayBlocks,
				"target-hash", finalPayload.ExecutionPayload.BlockHash,
				"target-num", uint64(finalPayload.ExecutionPayload.BlockNumber),
			)
			building = true
		} else {
			// TODO: we could probably do the building async,
			// to potentially stress out the builder
			if building {
				var timestamp uint64
				var nextPayload *eth.ExecutionPayloadEnvelope
				if payloadIndex > 1 {
					nextPayload = payloads[payloadIndex-2]
				} else {
					nextPayload = payloadAfterFinal
				}
				timestamp = uint64(nextPayload.ExecutionPayload.Timestamp)
				// give time to process
				time.Sleep(30 * time.Millisecond)
				info := eth.PayloadInfo{
					ID: *status.BuildingID,
					// use the timestamp from the actual next payload
					Timestamp: timestamp,
				}
				replayedPayload, err := co.buildClient.EngineAPIClient.GetPayload(ctx, info)
				if err != nil {
					co.logger.Error("error during build payload retrieval")
				}
				// validate the built payload, this should put the next head into the cache already,
				// so that the next call to this caused by syncUntil will be redundant, if
				// we didn't get a bad block
				replayStatus, err := co.newPayload(ctx, replayedPayload)
				s <- SyncStatus{
					Status:     *replayStatus,
					Error:      err,
					BuildingID: status.BuildingID,
					// this is the payload we got from the trusted client
					TargetPayload: nextPayload,
					BuiltPayload:  replayedPayload,
					Replay:        true,
				}
			}
		}
		payloadIndex--
	}
	return nil
}

func (co *ConsensusOrchestrator) forkChoice(
	ctx context.Context,
	engine
	block common.Hash,
	attrs *eth.PayloadAttributes) (*eth.PayloadStatusV1, *eth.PayloadID, error) {
	var err error

	fc := &eth.ForkchoiceState{
		HeadBlockHash: block,
		// TODO: do we need this properly set?
		// how does geth use this during block building?
		// SafeBlockHash:      common.Hash{},
		// FinalizedBlockHash: common.Hash{},
	}
	fcResult, err := co.validateClient.EngineAPIClient.ForkchoiceUpdate(ctx, fc, attrs)
	if err != nil {
		// TODO: could be unrecoverable, return wrapped error to exit on it
		return nil, nil, fmt.Errorf("failed to update forkchoice: %w", err)
	}
	if fcResult.PayloadStatus.Status == eth.ExecutionInvalid {
		err := errors.New(*fcResult.PayloadStatus.ValidationError)
		return &fcResult.PayloadStatus, fcResult.PayloadID, fmt.Errorf("%w: %w", errExecutionInvalid, err)
	}
	return &fcResult.PayloadStatus, fcResult.PayloadID, nil
}

func PayloadToBuildAttrs(payload eth.ExecutionPayloadEnvelope) *eth.PayloadAttributes {
	return &eth.PayloadAttributes{
		Timestamp:             payload.ExecutionPayload.Timestamp,
		PrevRandao:            payload.ExecutionPayload.PrevRandao,
		SuggestedFeeRecipient: payload.ExecutionPayload.FeeRecipient,
		Withdrawals:           payload.ExecutionPayload.Withdrawals,
		ParentBeaconBlockRoot: payload.ParentBeaconBlockRoot,
		// TODO: we don't know 100% if this is suited to
		// insert normal transactions into the payload execution like that.
		Transactions: payload.ExecutionPayload.Transactions,
		NoTxPool:     true,
		GasLimit:     &payload.ExecutionPayload.GasLimit,
	}
}
