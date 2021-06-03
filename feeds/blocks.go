package feeds

import (
	"context"
	"errors"
	"math/big"

	log "github.com/sirupsen/logrus"

	"fortify-node/domain"
	"fortify-node/ethereum"
	"fortify-node/utils"
)

var ErrEndBlockReached = errors.New("end block reached")

type BlockFeed interface {
	ForEachBlock(handler func(evt *domain.BlockEvent) error) error
}

type blockFeed struct {
	start   *big.Int
	end     *big.Int
	ctx     context.Context
	client  ethereum.Client
	cache   utils.Cache
	chainID *big.Int
	tracing bool
}

type BlockFeedConfig struct {
	Start   *big.Int
	End     *big.Int
	ChainID *big.Int
	Tracing bool
}

func (bf *blockFeed) initialize() error {
	if bf.start == nil {
		res, err := bf.client.BlockByNumber(bf.ctx, nil)
		if err != nil {
			return err
		}
		log.Debugf("retrieved block number %s", res.Number)

		bf.start, err = utils.HexToBigInt(res.Number)
		if err != nil {
			log.Errorf("error converting blocknum hex to bigint: %s", err.Error())
			return nil
		}
	}
	log.Infof("initialized block number %d", bf.start)

	if bf.chainID == nil {
		chainID, err := bf.client.ChainID(bf.ctx)
		if err != nil {
			return err
		}
		bf.chainID = chainID
	}
	log.Infof("initialized chainId %d", bf.chainID)

	return nil
}

func (bf *blockFeed) processReorg(parentHash string, handler func(evt *domain.BlockEvent) error) error {
	// don't process anything before start index
	currentHash := parentHash
	for {
		if bf.ctx.Err() != nil {
			log.Debug("processReorg, returning ctx err")
			return bf.ctx.Err()
		}
		if bf.cache.Exists(currentHash) {
			return nil
		}
		block, err := bf.client.BlockByHash(bf.ctx, currentHash)
		if err != nil {
			log.Errorf("reorg: err getting block: %s (skipping)", err.Error())
			return nil
		}
		blockNum, err := utils.HexToBigInt(block.Number)
		if err != nil {
			log.Errorf("error converting blocknum hex to bigint: %s", err.Error())
			return nil
		}

		var traces []domain.Trace
		if bf.tracing {
			traces, err = bf.client.TraceBlock(bf.ctx, blockNum)
			if err != nil {
				log.Errorf("error tracing block: %s", err.Error())
				return err
			}
		}

		if blockNum.Uint64() <= bf.start.Uint64() {
			// stop if prior to horizon
			return nil
		}
		evt := &domain.BlockEvent{EventType: domain.EventTypeReorg, Block: block, ChainID: bf.chainID, Traces: traces}
		if err := handler(evt); err != nil {
			return err
		}

		bf.cache.Add(currentHash)
		currentHash = block.ParentHash
	}
}

func (bf *blockFeed) ForEachBlock(handler func(evt *domain.BlockEvent) error) error {
	increment := big.NewInt(1)
	blockNum := big.NewInt(bf.start.Int64())

	for {
		if bf.ctx.Err() != nil {
			return bf.ctx.Err()
		}
		if bf.end != nil && blockNum.Uint64() > bf.end.Uint64() {
			return ErrEndBlockReached
		}

		var err error
		var traces []domain.Trace
		if bf.tracing {
			traces, err = bf.client.TraceBlock(bf.ctx, blockNum)
			if err != nil {
				log.Errorf("error tracing block: %s", err.Error())
				return err
			}
		}

		var block *domain.Block
		if len(traces) == 0 {
			block, err = bf.client.BlockByNumber(bf.ctx, blockNum)
		} else {
			// this forces the SAME block to be returned as traces (so that a re-org doesn't split it)
			hash := traces[0].BlockHash
			block, err = bf.client.BlockByHash(bf.ctx, *hash)
		}
		if err != nil {
			log.Errorf("error getting block: %s", err.Error())
			return err
		}

		evt := &domain.BlockEvent{EventType: domain.EventTypeBlock, Block: block, ChainID: bf.chainID, Traces: traces}
		if err := handler(evt); err != nil {
			return err
		}
		bf.cache.Add(block.Hash)
		if blockNum.Uint64() > bf.start.Uint64() {
			if err := bf.processReorg(block.ParentHash, handler); err != nil {
				log.Errorf("ForEachBlock: err from processReorg: %s", err.Error())
				return err
			}
		}
		blockNum.Add(blockNum, increment)
	}
}

func NewBlockFeed(ctx context.Context, client ethereum.Client, cfg BlockFeedConfig) (*blockFeed, error) {
	bf := &blockFeed{
		start:   cfg.Start,
		end:     cfg.End,
		ctx:     ctx,
		client:  client,
		cache:   utils.NewCache(10000),
		chainID: cfg.ChainID,
		tracing: cfg.Tracing,
	}
	if err := bf.initialize(); err != nil {
		return nil, err
	}
	return bf, nil
}
