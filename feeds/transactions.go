package feeds

import (
	"bytes"
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"OpenZeppelin/safe-node/clients"
	"OpenZeppelin/safe-node/protocol"
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
	blockCh   chan *BlockEvent
	txCh      chan *TransactionEvent
}

type TransactionEvent struct {
	EventType   EventType
	Block       *types.Block
	Transaction *types.Transaction
	Receipt     *types.Receipt
}

// ToMessage converts the TransactionEvent to the protocol.TransactionEvent message
func (t *TransactionEvent) ToMessage() (*protocol.TransactionEvent, error) {
	evtType := protocol.TransactionEvent_BLOCK
	if t.EventType == "reorg" {
		evtType = protocol.TransactionEvent_REORG
	}
	var tx protocol.TransactionEvent_EthTransaction
	var receipt protocol.TransactionEvent_EthReceipt
	um := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}

	if t.Transaction != nil {
		txJson, err := t.Transaction.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if err := um.Unmarshal(bytes.NewReader(txJson), &tx); err != nil {
			return nil, err
		}
	}

	if t.Receipt != nil {
		receiptJson, err := t.Receipt.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if err := um.Unmarshal(bytes.NewReader(receiptJson), &receipt); err != nil {
			return nil, err
		}
	}

	return &protocol.TransactionEvent{
		Type:        evtType,
		Transaction: &tx,
		Receipt:     &receipt,
	}, nil
}

func (tf *transactionFeed) streamBlocks() error {
	defer close(tf.blockCh)
	return tf.blockFeed.ForEachBlock(func(evt *BlockEvent) error {
		log.Debugf("block-iterator: blocks <- %d", evt.Block.NumberU64())
		tf.blockCh <- evt
		return nil
	})
}

func (tf *transactionFeed) streamTransactions() error {
	defer close(tf.txCh)
	for evt := range tf.blockCh {
		log.Debugf("tx-iterator: block(%d) processing", evt.Block.NumberU64())
		for _, tx := range evt.Block.Transactions() {
			select {
			case <-tf.ctx.Done():
				return tf.ctx.Err()
			default:
				if !tf.cache.ExistsAndAdd(tx.Hash().Hex()) {
					log.Debugf("tx-iterator: block(%d), txs <- %s", evt.Block.NumberU64(), tx.Hash().Hex())
					tf.txCh <- &TransactionEvent{EventType: evt.EventType, Block: evt.Block, Transaction: tx}
				}
			}
		}
	}
	return nil
}

func (tf *transactionFeed) getWorker(workerID int, handler func(evt *TransactionEvent) error) func() error {
	return func() error {
		for tx := range tf.txCh {
			log.Debugf("tx-processor(%d): block(%d) processing %s", workerID, tx.Block.NumberU64(), tx.Transaction.Hash().Hex())
			select {
			case <-tf.ctx.Done():
				log.Debugf("tx-processor(%d): context cancelled", workerID)
				return tf.ctx.Err()
			default:
				receipt, err := tf.client.TransactionReceipt(tf.ctx, tx.Transaction.Hash())
				if err != nil {
					log.Debugf("tx-processor(%d): block(%d) tx(%s) get receipt failed (skipping): %s", workerID, tx.Block.NumberU64(), tx.Transaction.Hash().Hex(), err.Error())
					continue
				}
				tx.Receipt = receipt
				log.Debugf("tx-processor(%d): block(%d) tx(%s) invoking handler", workerID, tx.Block.NumberU64(), tx.Transaction.Hash().Hex())
				if err := handler(tx); err != nil {
					log.Debugf("tx-processor(%d): block(%d) tx(%s) handler returned error, cancelling: %s", workerID, tx.Block.NumberU64(), tx.Transaction.Hash().Hex(), err.Error())
					return err
				}
			}
		}
		return nil
	}
}

// ForEachTransaction invokes a handler for each transactions on a network until cancelled or handler returns error
func (tf *transactionFeed) ForEachTransaction(handler func(evt *TransactionEvent) error) error {
	grp, _ := errgroup.WithContext(tf.ctx)

	// iterate over blocks
	grp.Go(tf.streamBlocks)

	// iterate over transactions, check for duplicates
	grp.Go(tf.streamTransactions)

	// because my tests weren't working and this was why
	if tf.workers < 1 {
		return errors.New("workers must be > 0")
	}

	// get receipt and invoke handler for each transaction (x workers)
	for i := 0; i < tf.workers; i++ {
		workerID := i
		grp.Go(tf.getWorker(workerID, handler))
	}

	// block until above all finish (when context is cancelled or error returns)
	return grp.Wait()
}

func NewTransactionFeed(ctx context.Context, client clients.EthClient, start *big.Int, workers int) (*transactionFeed, error) {
	blocks := make(chan *BlockEvent, 10)
	txs := make(chan *TransactionEvent, 100)
	blockFeed, err := NewBlockFeed(ctx, client, start)
	cache := utils.NewCache(1000000)
	if err != nil {
		return nil, err
	}
	return &transactionFeed{
		ctx: ctx, cache: cache, client: client, blockFeed: blockFeed, workers: workers, blockCh: blocks, txCh: txs,
	}, nil
}
