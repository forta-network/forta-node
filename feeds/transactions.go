package feeds

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"forta-network/forta-node/domain"
	"forta-network/forta-node/ethereum"
	"forta-network/forta-node/utils"
)

type transactionFeed struct {
	ctx       context.Context
	cache     utils.Cache
	client    ethereum.Client
	blockFeed BlockFeed
	workers   int
	blockCh   chan *domain.BlockEvent
	txCh      chan *domain.TransactionEvent
}

func (tf *transactionFeed) streamTransactions() error {
	defer close(tf.txCh)
	for {
		blockEvt, ok := <-tf.blockCh
		if !ok {
			return nil
		}

		log.Infof("tx-iterator: block(%s) processing %d transactions", blockEvt.Block.Number, len(blockEvt.Block.Transactions))
		for _, tx := range blockEvt.Block.Transactions {
			txTemp := tx
			select {
			case <-tf.ctx.Done():
				return tf.ctx.Err()
			default:
				if !tf.cache.ExistsAndAdd(tx.Hash) {
					log.Debugf("tx-iterator: block(%s), txs <- %s", blockEvt.Block.Number, tx.Hash)
					tf.txCh <- &domain.TransactionEvent{
						BlockEvt:    blockEvt,
						Transaction: &txTemp,
					}
				}
			}
		}
	}
}

func (tf *transactionFeed) getWorker(workerID int, handler func(evt *domain.TransactionEvent) error) func() error {
	return func() error {
		for tx := range tf.txCh {
			log.Debugf("tx-processor(%d): block(%s) processing %s", workerID, tx.BlockEvt.Block.Number, tx.Transaction.Hash)
			select {
			case <-tf.ctx.Done():
				log.Debugf("tx-processor(%d): context cancelled", workerID)
				return tf.ctx.Err()
			default:
				r, err := tf.client.TransactionReceipt(tf.ctx, tx.Transaction.Hash)
				if err != nil {
					log.Warnf("tx-processor(%d): block(%s) tx(%s) get receipt failed (skipping): %s", workerID, tx.BlockEvt.Block.Number, tx.Transaction.Hash, err.Error())
					continue
				}
				tx.Receipt = r

				if err := handler(tx); err != nil {
					log.Errorf("tx-processor(%d): block(%s) tx(%s) handler returned error, cancelling: %s", workerID, tx.BlockEvt.Block.Number, tx.Transaction.Hash, err.Error())
					return err
				}
			}
		}
		return nil
	}
}

// ForEachTransaction invokes a handler for each transactions on a network until cancelled or handler returns error
func (tf *transactionFeed) ForEachTransaction(blockHandler func(evt *domain.BlockEvent) error, txHandler func(evt *domain.TransactionEvent) error) error {
	grp, _ := errgroup.WithContext(tf.ctx)

	// iterate over blocks
	grp.Go(func() error {
		errCh := tf.blockFeed.Subscribe(func(evt *domain.BlockEvent) error {
			log.Debugf("block-iterator: blocks <- %s", evt.Block.Number)
			tf.blockCh <- evt
			var blockHandlerErr error
			if blockHandler != nil {
				blockHandlerErr = blockHandler(evt)
			}
			return blockHandlerErr
		})
		err := <-errCh
		close(tf.blockCh)
		if err == ErrEndBlockReached {
			return nil
		}
		return err
	})

	// iterate over transactions, check for duplicates
	grp.Go(tf.streamTransactions)

	// because my tests weren't working and this was why
	if tf.workers < 1 {
		return errors.New("workers must be > 0")
	}

	// get receipt and invoke handler for each transaction (x workers)
	for i := 0; i < tf.workers; i++ {
		workerID := i
		grp.Go(tf.getWorker(workerID, txHandler))
	}

	// block until above all finish (when context is cancelled or error returns)
	return grp.Wait()
}

func NewTransactionFeed(ctx context.Context, client ethereum.Client, blockFeed BlockFeed, workers int) (*transactionFeed, error) {
	blocks := make(chan *domain.BlockEvent, 10)
	txs := make(chan *domain.TransactionEvent, 100)
	cache := utils.NewCache(1000000)
	return &transactionFeed{
		ctx: ctx, cache: cache, client: client, blockFeed: blockFeed, workers: workers, blockCh: blocks, txCh: txs,
	}, nil
}
