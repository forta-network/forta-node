package feeds

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/domain"
	"github.com/forta-protocol/forta-node/ethereum"
	"github.com/forta-protocol/forta-node/utils"
)

var ErrEndBlockReached = errors.New("end block reached")

type bfHandler struct {
	Handler func(evt *domain.BlockEvent) error
	ErrCh   chan<- error
}

type blockFeed struct {
	start       *big.Int
	end         *big.Int
	ctx         context.Context
	client      ethereum.Client
	traceClient ethereum.Client
	handlers    []*bfHandler
	cache       utils.Cache
	chainID     *big.Int
	tracing     bool
	started     bool
	rateLimit   *time.Ticker
	maxBlockAge *time.Duration
}

type BlockFeedConfig struct {
	Start               *big.Int
	End                 *big.Int
	ChainID             *big.Int
	RateLimit           *time.Ticker
	Tracing             bool
	SkipBlocksOlderThan *time.Duration
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

		// should be a positive number
		if bf.start.Sign() <= 0 {
			return fmt.Errorf("got invalid block number during initialization: %d", bf.start.Uint64())
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

func (bf *blockFeed) IsStarted() bool {
	return bf.started
}

func (bf *blockFeed) Start() {
	if !bf.started {
		go bf.loop()
	}
}

//StartRange runs a specific set of blocks synchronously
func (bf *blockFeed) StartRange(start int64, end int64, rate int64) {
	if !bf.started {
		if rate > 0 {
			bf.rateLimit = time.NewTicker(time.Duration(rate) * time.Millisecond)
		}
		bf.start = big.NewInt(start)
		bf.end = big.NewInt(end)
		go bf.loop()
	}
}

func (bf *blockFeed) loop() {
	bf.started = true
	defer func() {
		bf.started = false
	}()
	err := bf.forEachBlock()
	if err == nil {
		return
	}
	if err != ErrEndBlockReached {
		log.Warnf("failed while processing blocks: %v", err)
	}
	for _, handler := range bf.handlers {
		handler.ErrCh <- err
	}
}

func (bf *blockFeed) Subscribe(handler func(evt *domain.BlockEvent) error) <-chan error {
	errCh := make(chan error)
	bf.handlers = append(bf.handlers, &bfHandler{
		Handler: handler,
		ErrCh:   errCh,
	})
	return errCh
}

func (bf *blockFeed) forEachBlock() error {
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
			traces, err = bf.traceClient.TraceBlock(bf.ctx, blockNum)
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
			continue
		}

		if err != nil {
			log.Errorf("error getting blocknumber: num=%s, %s", block.Number, err.Error())
			continue
		}
		logger := log.WithFields(log.Fields{
			"blockNum": blockNum.Uint64(),
			"blockHex": block.Number,
		})
		age, err := block.Age()
		if err != nil || age == nil {
			logger.Errorf("error getting age of block: ts=%s, %s", block.Timestamp, err.Error())
			continue
		}

		// if not too old
		if bf.maxBlockAge == nil || *age < *bf.maxBlockAge {
			evt := &domain.BlockEvent{EventType: domain.EventTypeBlock, Block: block, ChainID: bf.chainID, Traces: traces}
			for _, handler := range bf.handlers {
				if err := handler.Handler(evt); err != nil {
					return err
				}
			}
			bf.cache.Add(block.Hash)
		} else {
			logger.WithField("age", age).Warnf("ignoring block, older than %v", bf.maxBlockAge)
		}

		blockNum.Add(blockNum, increment)
		if bf.rateLimit != nil {
			<-bf.rateLimit.C
		}
	}
}

func NewBlockFeed(ctx context.Context, client ethereum.Client, traceClient ethereum.Client, cfg BlockFeedConfig) (*blockFeed, error) {
	bf := &blockFeed{
		start:       cfg.Start,
		end:         cfg.End,
		ctx:         ctx,
		client:      client,
		traceClient: traceClient,
		cache:       utils.NewCache(10000),
		chainID:     cfg.ChainID,
		tracing:     cfg.Tracing,
		rateLimit:   cfg.RateLimit,
	}
	if err := bf.initialize(); err != nil {
		return nil, err
	}
	return bf, nil
}
