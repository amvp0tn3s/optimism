package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
)

func (co *ConsensusOrchestrator) forkChoice(ctx context.Context, errChan chan<- error, block common.Hash) {
	var err error
	defer func() {
		errChan <- err
		close(errChan)
	}()

	payload, err := co.liveClient.PayloadByHash(ctx, block)
	if err != nil {
		err = fmt.Errorf("couldn't get replayed block payload: %w", err)
		return
	}
	co.logger.Info("got payload for block", "block", block, "payload", payload)

	fc := &eth.ForkchoiceState{
		HeadBlockHash: payload.ExecutionPayload.ParentHash,
		// TODO: do we need this properly set?
		// how does geth use this during block building?
		SafeBlockHash:      common.Hash{},
		FinalizedBlockHash: common.Hash{},
	}

	results, err := co.client.EngineAPIClient.ForkchoiceUpdate(ctx, fc, nil)
	if err != nil {
		err = fmt.Errorf("failed to update forkchoice: %w", err)
		return
	}
	co.logger.Info("forkchoice update result", "result", results)
}

func (co *ConsensusOrchestrator) replay(ctx context.Context, errChan chan<- error, block common.Hash) {
	var err error
	defer func() {
		errChan <- err
		close(errChan)
	}()

	payload, err := co.client.PayloadByHash(ctx, block)
	if err != nil {
		err = fmt.Errorf("couldn't get replayed block payload: %w", err)
		return
	}

	// reorg to the parent block
	fc := &eth.ForkchoiceState{
		HeadBlockHash: payload.ExecutionPayload.ParentHash,
		// TODO: do we need this properly set?
		// how does geth use this during block building?
		SafeBlockHash:      common.Hash{},
		FinalizedBlockHash: common.Hash{},
	}
	// and supply the same attributes again for the block we want to have rebuild.
	// We (potentially) can directly pass the transactions and bypass the mempool
	// with the op-modifications of the engineAPI
	attributes := &eth.PayloadAttributes{
		Timestamp:             payload.ExecutionPayload.Timestamp,
		PrevRandao:            payload.ExecutionPayload.PrevRandao,
		SuggestedFeeRecipient: payload.ExecutionPayload.FeeRecipient,
		Withdrawals:           payload.ExecutionPayload.Withdrawals,
		ParentBeaconBlockRoot: payload.ParentBeaconBlockRoot,
		// TODO: we don't know 100% if this is suited to
		// insert normal transactions into the payload execution like that.
		Transactions: payload.ExecutionPayload.Transactions,
		NoTxPool:     true,
		// TODO: is this required to be set?
		GasLimit: &payload.ExecutionPayload.GasLimit,
	}

	results, err := co.client.EngineAPIClient.ForkchoiceUpdate(ctx, fc, attributes)
	if err != nil {
		err = fmt.Errorf("failed to update forkchoice: %w", err)
		return
	}
	time.Sleep(time.Second)

	info := eth.PayloadInfo{
		ID:        *results.PayloadID,
		Timestamp: uint64(attributes.Timestamp + 1),
	}
	// XXX: do we have to do another update in order to seal the block?
	replayedPayload, err := co.client.EngineAPIClient.GetPayload(ctx, info)
	// TODO: do some sanity checks in comparison with the old payload
	_ = replayedPayload

	// trigger the validation path:
	validationStatus, err := co.client.EngineAPIClient.NewPayload(ctx, payload.ExecutionPayload, payload.ParentBeaconBlockRoot)
	_ = validationStatus

	// TODO: log:
	// validationStatus.Status
	// validationStatus.LatestValidHash
	// validationStatus.ValidationError
}
