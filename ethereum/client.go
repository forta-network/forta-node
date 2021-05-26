package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"

	"OpenZeppelin/fortify-node/domain"
	"OpenZeppelin/fortify-node/utils"
)

type rpcClient interface {
	Close()
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

// Client is an interface encompassing all ethereum actions
type Client interface {
	Close()
	BlockByHash(ctx context.Context, hash string) (*domain.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*domain.Block, error)
	BlockNumber(ctx context.Context) (*big.Int, error)
	TransactionReceipt(ctx context.Context, txHash string) (*domain.TransactionReceipt, error)
	ChainID(ctx context.Context) (*big.Int, error)
	TraceBlock(ctx context.Context, number *big.Int) ([]domain.Trace, error)
}

const blocksByNumber = "eth_getBlockByNumber"
const blocksByHash = "eth_getBlockByHash"
const blockNumber = "eth_blockNumber"
const transactionReceipt = "eth_getTransactionReceipt"
const traceBlock = "trace_block"
const chainId = "eth_chainId"

var permanentErrors = []string{"method not found"}

var minBackoff = 1 * time.Second
var maxBackoff = 1 * time.Minute

// streamEthClient wraps a go-ethereum client purpose-built for streaming txs (with long retries/timeouts)
type streamEthClient struct {
	rpcClient rpcClient
}

type RetryOptions struct {
	MaxElapsedTime *time.Duration
	MinBackoff     *time.Duration
	MaxBackoff     *time.Duration
}

// Close invokes close on the underlying client
func (e streamEthClient) Close() {
	e.rpcClient.Close()
}

func isPermanentError(err error) bool {
	if err == nil {
		return false
	}
	for _, pe := range permanentErrors {
		if strings.Contains(strings.ToLower(err.Error()), pe) {
			return true
		}
	}
	return false
}

// withBackoff wraps an operation in an exponential backoff logic
func withBackoff(ctx context.Context, name string, operation func(ctx context.Context) error, options RetryOptions) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = maxBackoff
	bo.InitialInterval = minBackoff
	if options.MinBackoff != nil {
		bo.InitialInterval = *options.MinBackoff
	}
	if options.MaxBackoff != nil {
		bo.MaxInterval = *options.MaxBackoff
	}
	if options.MaxElapsedTime != nil {
		bo.MaxElapsedTime = *options.MaxElapsedTime
	}
	err := backoff.Retry(func() error {
		if ctx.Err() != nil {
			return backoff.Permanent(ctx.Err())
		}
		tCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		defer cancel()
		err := operation(tCtx)

		//any non-retriable failure errors can be listed here
		if ctx.Err() != nil {
			log.Debugf("%s context cancelled", name)
			return backoff.Permanent(ctx.Err())
		}
		if err != nil && isPermanentError(err) {
			log.Errorf("backoff permanent error: %s", err.Error())
			return backoff.Permanent(err)
		} else if err != nil {
			log.Debugf("%s failed...retrying: %s", name, err.Error())
		}
		return err
	}, bo)
	if err != nil {
		log.Errorf("%s failed with error: %s", name, err.Error())
	}
	return err
}

func pointDur(d time.Duration) *time.Duration {
	return &d
}

// BlockByHash returns the block by hash
func (e streamEthClient) BlockByHash(ctx context.Context, hash string) (*domain.Block, error) {
	name := fmt.Sprintf("%s(%s)", blocksByHash, hash)
	log.Debugf(name)
	var result domain.Block
	err := withBackoff(ctx, name, func(ctx context.Context) error {
		return e.rpcClient.CallContext(ctx, &result, blocksByHash, hash, true)
	}, RetryOptions{
		MaxElapsedTime: pointDur(12 * time.Hour),
		MaxBackoff:     pointDur(15 * time.Second),
	})
	return &result, err
}

// TraceBlock returns the traced block
func (e streamEthClient) TraceBlock(ctx context.Context, number *big.Int) ([]domain.Trace, error) {
	name := fmt.Sprintf("%s(%s)", traceBlock, number)
	log.Debugf(name)
	var result []domain.Trace
	err := withBackoff(ctx, name, func(ctx context.Context) error {
		return e.rpcClient.CallContext(ctx, &result, traceBlock, utils.BigIntToHex(number))
	}, RetryOptions{
		MaxElapsedTime: pointDur(12 * time.Hour),
		MaxBackoff:     pointDur(15 * time.Second),
	})
	return result, err
}

// BlockByNumber returns the block by number
func (e streamEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*domain.Block, error) {
	var result domain.Block
	num := "latest"
	if number != nil {
		num = utils.BigIntToHex(number)
	}
	name := fmt.Sprintf("%s(%s)", blocksByNumber, num)
	log.Debugf(name)

	err := withBackoff(ctx, name, func(ctx context.Context) error {
		return e.rpcClient.CallContext(ctx, &result, blocksByNumber, num, true)
	}, RetryOptions{
		MaxElapsedTime: pointDur(12 * time.Hour),
		MaxBackoff:     pointDur(15 * time.Second),
	})
	return &result, err
}

// BlockNumber returns the latest block number
func (e streamEthClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	log.Debugf(blockNumber)
	var result string
	err := withBackoff(ctx, blockNumber, func(ctx context.Context) error {
		return e.rpcClient.CallContext(ctx, &result, blockNumber)
	}, RetryOptions{
		MaxElapsedTime: pointDur(12 * time.Hour),
	})
	if err != nil {
		return nil, err
	}
	return utils.HexToBigInt(result)
}

// ChainID gets the chainID for a network
func (e streamEthClient) ChainID(ctx context.Context) (*big.Int, error) {
	log.Debugf(chainId)
	var result string
	err := withBackoff(ctx, chainId, func(ctx context.Context) error {
		return e.rpcClient.CallContext(ctx, &result, chainId)
	}, RetryOptions{
		MaxElapsedTime: pointDur(1 * time.Minute),
	})
	if err != nil {
		return nil, err
	}
	return utils.HexToBigInt(result)
}

// TransactionReceipt returns the receipt for a transaction
func (e streamEthClient) TransactionReceipt(ctx context.Context, txHash string) (*domain.TransactionReceipt, error) {
	name := fmt.Sprintf("%s(%s)", transactionReceipt, txHash)
	log.Debugf(name)
	var result domain.TransactionReceipt
	err := withBackoff(ctx, name, func(ctx context.Context) error {
		return e.rpcClient.CallContext(ctx, &result, transactionReceipt, txHash)
	}, RetryOptions{
		MaxElapsedTime: pointDur(5 * time.Minute),
	})
	return &result, err
}

// NewStreamEthClient creates a new ethereum client
func NewStreamEthClient(ctx context.Context, url string) (*streamEthClient, error) {
	//TODO: consider NewClient with a custom RPC so that one can inject headers
	rpcClient, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return &streamEthClient{rpcClient: rpcClient}, nil
}
