package feeds

import (
	"context"
	"math/big"

	"github.com/forta-network/forta-node/domain"

	"github.com/forta-network/forta-node/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	eth "github.com/forta-network/forta-node/ethereum"
)

type logFeed struct {
	ctx        context.Context
	url        string
	startBlock *big.Int
	endBlock   *big.Int
	topics     [][]string
	addresses  []string
	client     eth.Client

	blockCh chan *domain.Block
}

func (l *logFeed) ForEachLog(handler func(logEntry types.Log) error) error {
	eg, ctx := errgroup.WithContext(l.ctx)

	addrs := make([]common.Address, 0, len(l.addresses))
	for _, addr := range l.addresses {
		addrs = append(addrs, common.HexToAddress(addr))
	}

	var topics [][]common.Hash
	for _, topicSet := range l.topics {
		var topicPosition []common.Hash
		for _, topic := range topicSet {
			topicHash := common.HexToHash(topic)
			topicPosition = append(topicPosition, topicHash)
		}
		topics = append(topics, topicPosition)
	}

	currentBlock := l.startBlock
	increment := big.NewInt(1)
	eg.Go(func() error {
		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			blk, err := l.client.BlockByNumber(l.ctx, currentBlock)
			if err != nil {
				log.Error("error while getting latest block number:", err)
				return err
			}

			// initialize current block if nil
			if currentBlock == nil {
				currentBlock, err = utils.HexToBigInt(blk.Number)
				if err != nil {
					log.Errorf("error while converting latest block number: %s, %s", blk.Number, err)
					return err
				}
			}

			q := ethereum.FilterQuery{
				FromBlock: currentBlock,
				ToBlock:   currentBlock,
				Addresses: addrs,
				Topics:    topics,
			}
			logs, err := l.client.GetLogs(l.ctx, q)
			if err != nil {
				return err
			}
			for _, lg := range logs {
				if err := handler(lg); err != nil {
					log.Error("handler returned error, exiting log subscription:", err)
					return err
				}
			}
			if l.blockCh != nil {
				l.blockCh <- blk
			}

			currentBlock = currentBlock.Add(currentBlock, increment)
		}
		return nil
	})
	log.Infof("subscribed to logs: address=%v, topics=%v, startBlock=%s, endBlock=%s", l.addresses, l.topics, l.startBlock, l.endBlock)
	defer func() {
		log.Warn("log subscription closed")
	}()
	return eg.Wait()
}

type LogFeedConfig struct {
	Topics    [][]string
	Addresses []string
}

func NewLogFeed(ctx context.Context, client eth.Client, blockCh chan *domain.Block, cfg LogFeedConfig) (*logFeed, error) {
	return &logFeed{
		ctx:       ctx,
		client:    client,
		topics:    cfg.Topics,
		addresses: cfg.Addresses,
		blockCh:   blockCh,
	}, nil
}
