package feeds

import (
	"context"
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"OpenZeppelin/fortify-node/domain"
	"OpenZeppelin/fortify-node/ethereum"
	"OpenZeppelin/fortify-node/utils"
)

type TransactionFeed interface {
	ForEachTransaction(blockHandler func(evt *domain.BlockEvent) error, txHandler func(evt *domain.TransactionEvent) error) error
}

type transactionFeed struct {
	ctx       context.Context
	cache     utils.Cache
	client    ethereum.Client
	blockFeed BlockFeed
	workers   int
	blockCh   chan *domain.BlockEvent
	blocksOut chan *domain.BlockEvent
	txCh      chan *domain.TransactionEvent
}

func (tf *transactionFeed) streamBlocks() error {
	defer close(tf.blockCh)
	return tf.blockFeed.ForEachBlock(func(evt *domain.BlockEvent) error {
		log.Debugf("block-iterator: blocks <- %s", evt.Block.Number)
		tf.blockCh <- evt
		return nil
	})
}

func (tf *transactionFeed) streamTransactions() error {
	defer close(tf.txCh)
	for evt := range tf.blockCh {
		log.Debugf("tx-iterator: block(%s) processing", evt.Block.Number)
		for _, tx := range evt.Block.Transactions {
			select {
			case <-tf.ctx.Done():
				return tf.ctx.Err()
			default:
				if !tf.cache.ExistsAndAdd(tx.Hash) {
					log.Debugf("tx-iterator: block(%s), txs <- %s", evt.Block.Number, tx.Hash)
					tf.txCh <- &domain.TransactionEvent{
						BlockEvt:    evt,
						Transaction: &tx,
					}
				}
			}
		}
	}
	return nil
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
				receipt, err := tf.client.TransactionReceipt(tf.ctx, tx.Transaction.Hash)
				if err != nil {
					log.Debugf("tx-processor(%d): block(%s) tx(%s) get receipt failed (skipping): %s", workerID, tx.BlockEvt.Block.Number, tx.Transaction.Hash, err.Error())
					continue
				}
				tx.Receipt = receipt
				log.Debugf("tx-processor(%d): block(%s) tx(%s) invoking handler", workerID, tx.BlockEvt.Block.Number, tx.Transaction.Hash)
				if err := handler(tx); err != nil {
					log.Debugf("tx-processor(%d): block(%s) tx(%s) handler returned error, cancelling: %s", workerID, tx.BlockEvt.Block.Number, tx.Transaction.Hash, err.Error())
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
		defer close(tf.blockCh)
		return tf.blockFeed.ForEachBlock(func(evt *domain.BlockEvent) error {
			log.Debugf("block-iterator: blocks <- %s", evt.Block.Number)
			tf.blockCh <- evt
			return blockHandler(evt)
		})
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

func NewTransactionFeed(ctx context.Context, client ethereum.Client, chainID *big.Int, start *big.Int, workers int) (*transactionFeed, error) {
	blocks := make(chan *domain.BlockEvent, 10)
	blocksOut := make(chan *domain.BlockEvent)
	txs := make(chan *domain.TransactionEvent, 100)
	blockFeed, err := NewBlockFeed(ctx, client, chainID, start)
	cache := utils.NewCache(1000000)
	if err != nil {
		return nil, err
	}
	return &transactionFeed{
		ctx: ctx, cache: cache, client: client, blockFeed: blockFeed, workers: workers, blockCh: blocks, blocksOut: blocksOut, txCh: txs,
	}, nil
}
