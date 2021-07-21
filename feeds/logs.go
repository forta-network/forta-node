package feeds

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

type EthClient interface {
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	Close()
}

type logFeed struct {
	ctx    context.Context
	client EthClient
	query  ethereum.FilterQuery
}

func (l *logFeed) ForEachLog(handler func(log types.Log) error) error {
	logs := make(chan types.Log)
	sub, err := l.client.SubscribeFilterLogs(l.ctx, l.query, logs)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(l.ctx)
	ticker := time.NewTicker(100 * time.Millisecond)

	eg.Go(func() error {
		defer l.client.Close()
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-sub.Err():
				return err
			}
		}
	})

	eg.Go(func() error {
		for ethLog := range logs {
			if err := handler(ethLog); err != nil {
				return err
			}
		}
		return nil
	})

	return eg.Wait()
}

func NewLogFeed(ctx context.Context, wssUrl string, contractAddrs []string) (*logFeed, error) {
	client, err := ethclient.Dial(wssUrl)
	if err != nil {
		return nil, err
	}

	addrs := make([]common.Address, 0, len(contractAddrs))
	for _, addr := range contractAddrs {
		addrs = append(addrs, common.HexToAddress(addr))
	}

	return &logFeed{
		ctx:    ctx,
		client: client,
		query: ethereum.FilterQuery{
			Addresses: addrs,
		},
	}, nil
}
