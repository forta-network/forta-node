package feeds

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/domain"
	"OpenZeppelin/fortify-node/ethereum"
	"OpenZeppelin/fortify-node/utils"
)

type BlockFeed interface {
	ForEachBlock(handler func(evt *domain.BlockEvent) error) error
}

type blockFeed struct {
	start   *big.Int
	ctx     context.Context
	client  ethereum.Client
	cache   utils.Cache
	chainID *big.Int
}

func (bf *blockFeed) initialize() error {
	if bf.start == nil {
		res, err := bf.client.BlockByNumber(bf.ctx, nil)
		if err != nil {
			return err
		}
		log.Debugf("retrieved block number %d", res.Number())
		bf.start = big.NewInt(int64(res.NumberU64()))
	}
	log.Infof("initialized block number %d", bf.start)

	chainID, err := bf.client.ChainID(bf.ctx)
	if err != nil {
		return err
	}
	bf.chainID = chainID
	log.Infof("initialized chain id %d", bf.chainID)
	return nil
}

func (bf *blockFeed) processReorg(parentHash common.Hash, handler func(evt *domain.BlockEvent) error) error {
	// don't process anything before start index
	currentHash := parentHash
	for {
		if bf.ctx.Err() != nil {
			log.Debug("processReorg, returning ctx err")
			return bf.ctx.Err()
		}
		if bf.cache.Exists(currentHash.Hex()) {
			return nil
		}
		block, err := bf.client.BlockByHash(bf.ctx, currentHash)
		if err != nil {
			log.Errorf("reorg: err getting block: %s (skipping)", err.Error())
			return nil
		}

		traces, err := bf.client.TraceBlock(bf.ctx, block.Number())
		if err != nil {
			log.Errorf("reorg: error tracing block: %s (skipping)", err.Error())
			return nil
		}

		if block.NumberU64() <= bf.start.Uint64() {
			// stop if prior to horizon
			return nil
		}
		evt := &domain.BlockEvent{EventType: domain.EventTypeReorg, Block: block, ChainID: bf.chainID, Traces: traces}
		if err := handler(evt); err != nil {
			return err
		}

		bf.cache.Add(currentHash.Hex())
		currentHash = block.ParentHash()
	}
}

func (bf *blockFeed) ForEachBlock(handler func(evt *domain.BlockEvent) error) error {
	increment := big.NewInt(1)
	blockNum := big.NewInt(bf.start.Int64())

	for {
		if bf.ctx.Err() != nil {
			return bf.ctx.Err()
		}

		traces, err := bf.client.TraceBlock(bf.ctx, blockNum)
		if err != nil {
			log.Errorf("error tracing block: %s", err.Error())
			return err
		}

		var block *types.Block
		if len(traces) == 0 {
			block, err = bf.client.BlockByNumber(bf.ctx, blockNum)
		} else {
			// this forces the SAME block to be returned as traces (so that a re-org doesn't split it)
			hash := traces[0].BlockHash
			block, err = bf.client.BlockByHash(bf.ctx, common.HexToHash(*hash))
		}
		if err != nil {
			log.Errorf("error getting block: %s", err.Error())
			return err
		}

		evt := &domain.BlockEvent{EventType: domain.EventTypeBlock, Block: block, ChainID: bf.chainID, Traces: traces}
		if err := handler(evt); err != nil {
			return err
		}
		bf.cache.Add(block.Hash().Hex())
		if blockNum.Uint64() > bf.start.Uint64() {
			if err := bf.processReorg(block.ParentHash(), handler); err != nil {
				log.Errorf("ForEachBlock: err from processReorg: %s", err.Error())
				return err
			}
		}
		blockNum.Add(blockNum, increment)
	}
}

func NewBlockFeed(ctx context.Context, client ethereum.Client, start *big.Int) (*blockFeed, error) {
	bf := &blockFeed{
		start: start, ctx: ctx, client: client, cache: utils.NewCache(10000),
	}
	if err := bf.initialize(); err != nil {
		return nil, err
	}
	return bf, nil
}
