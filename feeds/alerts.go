package feeds

import (
	"context"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/forta-protocol/forta-node/domain"

	"github.com/forta-protocol/forta-node/contracts"

	"github.com/ethereum/go-ethereum/core/types"
)

const AlertBatchSignature = "AlertBatch(bytes32,address,uint256,uint256,uint256,uint256,uint256,string)"

//AlertBatchTopic is the topic value for the AlertBatch event, which can be used for filtering
var AlertBatchTopic = crypto.Keccak256Hash([]byte(AlertBatchSignature)).Hex()

type alertFeed struct {
	ctx context.Context
	lf  LogFeed
}

//ForEachAlert wraps a LogFeed.ForEachLog invocation and parses out the alert object
func (af *alertFeed) ForEachAlert(handler func(blk *domain.Block, batch *contracts.AlertsAlertBatch) error, finishBlockHandler func(blk *domain.Block) error) error {

	// cache by address so we don't over-allocate
	filterers := make(map[string]*contracts.AlertsFilterer)
	return af.lf.ForEachLog(func(blk *domain.Block, logEntry types.Log) error {
		if af.ctx.Err() != nil {
			return af.ctx.Err()
		}

		// filterers are per-contract address, this cache prevents overallocation
		if _, ok := filterers[logEntry.Address.Hex()]; !ok {
			f, err := contracts.NewAlertsFilterer(logEntry.Address, nil)
			if err != nil {
				return err
			}
			filterers[logEntry.Address.Hex()] = f
		}

		batch, err := filterers[logEntry.Address.Hex()].ParseAlertBatch(logEntry)
		if err != nil {
			return err
		}
		if batch != nil {
			return handler(blk, batch)
		}
		return nil
	}, finishBlockHandler)
}

// NewAlertFeed creates a new alert feed from a logFeed
func NewAlertFeed(ctx context.Context, lf LogFeed) (*alertFeed, error) {
	return &alertFeed{
		ctx: ctx,
		lf:  lf,
	}, nil
}
