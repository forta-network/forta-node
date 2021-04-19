package feeds

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"OpenZeppelin/safe-node/clients"
	"OpenZeppelin/safe-node/utils"
)

type TransactionFeed interface {
	ForEachTransaction(handler func(evt *TransactionEvent) error) error
}

type transactionFeed struct {
	ctx       context.Context
	cache     utils.Cache
	client    clients.EthClient
	blockFeed BlockFeed
	workers   int
}

type TransactionEvent struct {
	BlockEvent  *BlockEvent
	Transaction *types.Transaction
	Receipt     *types.Receipt
}

func (tf *transactionFeed) ForEachTransaction(handler func(evt *TransactionEvent) error) error {
	blocks := make(chan *BlockEvent, 10)
	txs := make(chan *TransactionEvent, 100)

	grp, _ := errgroup.WithContext(tf.ctx)

	// iterate over blocks
	grp.Go(func() error {
		defer close(blocks)
		return tf.blockFeed.ForEachBlock(func(evt *BlockEvent) error {
			log.Debugf("block-iterator: blocks <- %d", evt.Block.NumberU64())
			blocks <- evt
			return nil
		})
	})

	// iterate over transactions, check for duplicates
	grp.Go(func() error {
		defer close(txs)
		for evt := range blocks {
			log.Debugf("tx-iterator: block(%d) processing", evt.Block.NumberU64())
			for _, tx := range evt.Block.Transactions() {
				select {
				case <-tf.ctx.Done():
					return tf.ctx.Err()
				default:
					if !tf.cache.ExistsAndAdd(tx.Hash().Hex()) {
						log.Debugf("tx-iterator: block(%d), txs <- %s", evt.Block.NumberU64(), tx.Hash().Hex())
						txs <- &TransactionEvent{BlockEvent: evt, Transaction: tx}
					}
				}
			}
		}
		return nil
	})

	// get receipt and invoke handler for each transaction (x workers)
	for i := 0; i < tf.workers; i++ {
		workerID := i
		grp.Go(func() error {
			for tx := range txs {
				log.Debugf("tx-processor(%d): block(%d) processing %s", workerID, tx.BlockEvent.Block.NumberU64(), tx.Transaction.Hash().Hex())
				select {
				case <-tf.ctx.Done():
					log.Debugf("tx-processor(%d): context cancelled", workerID)
					return tf.ctx.Err()
				default:
					receipt, err := tf.client.TransactionReceipt(tf.ctx, tx.Transaction.Hash())
					if err != nil {
						log.Debugf("tx-processor(%d): block(%d) tx(%s) get receipt failed (skipping): %s", workerID, tx.BlockEvent.Block.NumberU64(), tx.Transaction.Hash().Hex(), err.Error())
						continue
					}
					tx.Receipt = receipt
					log.Debugf("tx-processor(%d): block(%d) tx(%s) invoking handler", workerID, tx.BlockEvent.Block.NumberU64(), tx.Transaction.Hash().Hex())
					if err := handler(tx); err != nil {
						log.Debugf("tx-processor(%d): block(%d) tx(%s) handler returned error, cancelling: %s", workerID, tx.BlockEvent.Block.NumberU64(), tx.Transaction.Hash().Hex(), err.Error())
						return err
					}
				}
			}
			return nil
		})
	}

	return grp.Wait()
}

func NewTransactionFeed(ctx context.Context, client clients.EthClient, start *big.Int) (*transactionFeed, error) {
	blockFeed, err := NewBlockFeed(ctx, client, start)
	if err != nil {
		return nil, err
	}
	return &transactionFeed{
		ctx, utils.NewCache(1000000), client, blockFeed, 10,
	}, nil
}
