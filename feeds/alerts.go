package feeds

import (
	"context"

	"github.com/forta-network/forta-node/domain"

	"github.com/forta-network/forta-node/contracts"

	"github.com/ethereum/go-ethereum/core/types"
)

//AlertBatchTopic is the topic value for the AlertBatch event, which can be used for filtering
const AlertBatchTopic = "0x36cde681f44e056b0e848fa24ffca3217ac9323460feeacf1a8ad8da28daf924"

type alertFeed struct {
	ctx context.Context
	lf  LogFeed
}

//ForEachAlert wraps a LogFeed.ForEachLog invocation and parses out the alert object
func (af *alertFeed) ForEachAlert(blockHandler func(blk *domain.Block) error, handler func(batch *contracts.AlertsAlertBatch) error) error {

	// cache by address so we don't over-allocate
	filterers := make(map[string]*contracts.AlertsFilterer)
	return af.lf.ForEachLog(blockHandler, func(logEntry types.Log) error {
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
			return handler(batch)
		}
		return nil
	})
}

// NewAlertFeed creates a new alert feed from a logFeed
func NewAlertFeed(ctx context.Context, lf LogFeed) (*alertFeed, error) {
	return &alertFeed{
		ctx: ctx,
		lf:  lf,
	}, nil
}
