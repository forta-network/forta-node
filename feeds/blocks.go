package feeds

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/clients/health"
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

	lastBlockByNumberReq health.TimeTracker
	lastBlockByNumberErr health.ErrorTracker
	lastTraceBlockReq    health.TimeTracker
	lastTraceBlockErr    health.ErrorTracker
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
	if err := bf.initialize(); err != nil {
		log.WithError(err).Panic("failed to initialize")
	}

	bf.started = true
	defer func() {
		bf.started = false
	}()
	err := bf.forEachBlock()
	if err == nil {
		return
	}
	if err != ErrEndBlockReached {
		log.WithError(err).Warn("failed while processing blocks")
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
		if bf.rateLimit != nil {
			<-bf.rateLimit.C
		}

		var err error
		var traces []domain.Trace
		if bf.tracing {
			traces, err = bf.traceClient.TraceBlock(bf.ctx, blockNum)
			if err != nil {
				log.WithError(err).Error("error tracing block")
			}
		}

		block, err := bf.client.BlockByNumber(bf.ctx, blockNum)
		if err != nil {
			log.WithError(err).Error("error getting block")
			continue
		}

		if len(traces) > 0 && block.Hash != utils.String(traces[0].BlockHash) {
			log.WithFields(log.Fields{
				"ethereumBlockHash": block.Hash,
				"traceBlockHash":    utils.String(traces[0].BlockHash),
			}).Warn("mismatching block hash from ethereum and trace apis - ignoring traces")
			traces = nil
		}

		if err != nil {
			log.WithError(err).Errorf("error getting blocknumber: num=%s", block.Number)
			continue
		}
		logger := log.WithFields(log.Fields{
			"blockNum": blockNum.Uint64(),
			"blockHex": block.Number,
		})
		// if not too old
		tooOld, age := blockIsTooOld(block, bf.maxBlockAge)
		if !tooOld {
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
	}
}

func blockIsTooOld(block *domain.Block, maxAge *time.Duration) (bool, *time.Duration) {
	if maxAge == nil {
		return false, nil
	}
	age, err := block.Age()
	if err != nil || age == nil {
		log.WithFields(log.Fields{
			"blockHex": block.Number,
		}).WithError(err).Errorf("error getting age of block")
		return false, age
	}
	return *age > *maxAge, age
}

// Name returns the name of this implementation.
func (bf *blockFeed) Name() string {
	return "block-feed"
}

// Health implements the health.Reporter interface.
func (bf *blockFeed) Health() health.Reports {
	return health.Reports{
		bf.lastBlockByNumberReq.GetReport("event.block-by-number.time"),
		bf.lastBlockByNumberErr.GetReport("event.block-by-number.error"),
		bf.lastBlockByNumberReq.GetReport("event.block-by-number.time"),
		bf.lastBlockByNumberErr.GetReport("event.block-by-number.error"),
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
		maxBlockAge: cfg.SkipBlocksOlderThan,
	}
	return bf, nil
}
