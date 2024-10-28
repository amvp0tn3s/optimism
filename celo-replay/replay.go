package replay

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"

	opnode "github.com/ethereum-optimism/optimism/op-node"
	"github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	L2     node.L2EndpointSetup
	Rollup node.Config
}

type Replayer struct {
	client *sources.EngineClient

	Block common.Hash
}

func Setup(ctx context.Context, cliCtx *cli.Context) (*Replayer, error) {
	var log log.Logger

	rollupCfg, err := opnode.NewRollupConfigFromCLI(log, cliCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to build rollup config: %w", err)
	}
	l2Cfg, err := opnode.NewL2EndpointConfig(cliCtx, log)
	if err != nil {
		return nil, fmt.Errorf("failed to load l2 endpoints info: %w", err)
	}
	rpcClient, rpcCfg, err := l2Cfg.Setup(ctx, log, rollupCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup L2 execution-engine RPC client: %w", err)
	}

	engineClient, err := sources.NewEngineClient(
		rpcClient, log, nil, rpcCfg,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate engine client: %w", err)
	}
	return &Replayer{
		client: engineClient,
		// TODO: parse in the replay block
		Block: common.Hash{},
	}, nil
}

func (r *Replayer) Start(ctx context.Context) error {
	return nil
}
func (r *Replayer) Stop(ctx context.Context) error {
	return nil
}
func (r *Replayer) Stopped() bool {
	return false
}

func (r *Replayer) replay(ctx context.Context) error {
	payload, err := r.client.PayloadByHash(ctx, r.Block)
	if err != nil {
		return fmt.Errorf("couldn't get replayed block payload: %w", err)
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

	results, err := r.client.EngineAPIClient.ForkchoiceUpdate(ctx, fc, attributes)
	if err != nil {
		return fmt.Errorf("failed to update forkchoice: %w", err)
	}
	time.Sleep(time.Second)

	info := eth.PayloadInfo{
		ID:        *results.PayloadID,
		Timestamp: uint64(attributes.Timestamp + 1),
	}
	// XXX: do we have to do another update in order to seal the block?
	replayedPayload, err := r.client.EngineAPIClient.GetPayload(ctx, info)
	// TODO: do some sanity checks in comparison with the old payload
	_ = replayedPayload

	// trigger the validation path:
	validationStatus, err := r.client.EngineAPIClient.NewPayload(ctx, payload.ExecutionPayload, payload.ParentBeaconBlockRoot)
	_ = validationStatus

	// TODO: log:
	// validationStatus.Status
	// validationStatus.LatestValidHash
	// validationStatus.ValidationError

	return nil
}
