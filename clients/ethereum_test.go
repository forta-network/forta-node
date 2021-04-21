package clients

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mocks "OpenZeppelin/safe-node/clients/mocks"
	"OpenZeppelin/safe-node/testutils"
)

const testBlockHash = "0x4fc0862e76691f5312964883954d5c2db35e2b8f7a4f191775a4f50c69804a8d"
const testTxHash = "0x9b9cc76d6b3b51976b1396a5b417b3bf3f4b39b8fe080e4a5aef39d02be882df"

var testErr = errors.New("test err")

func initClient(t *testing.T) (*ethClient, *mocks.MockEthClient, context.Context) {
	minBackoff = 1 * time.Millisecond
	maxBackoff = 1 * time.Millisecond
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mocks.NewMockEthClient(ctrl)

	return &ethClient{client}, client, ctx
}

func TestEthClient_BlockByHash(t *testing.T) {
	ethClient, client, ctx := initClient(t)

	hash := common.HexToHash(testBlockHash)
	block := testutils.TestBlock(1)

	// verify retry
	client.EXPECT().BlockByHash(gomock.Any(), hash).Return(nil, testErr).Times(1)
	client.EXPECT().BlockByHash(gomock.Any(), hash).Return(block, nil).Times(1)

	res, err := ethClient.BlockByHash(ctx, hash)
	assert.NoError(t, err)
	assert.Equal(t, block, res)
}

func TestEthClient_BlockByNumber(t *testing.T) {
	ethClient, client, ctx := initClient(t)

	block := testutils.TestBlock(1)
	num := big.NewInt(1)

	// verify retry
	client.EXPECT().BlockByNumber(gomock.Any(), num).Return(nil, testErr).Times(1)
	client.EXPECT().BlockByNumber(gomock.Any(), num).Return(block, nil).Times(1)

	res, err := ethClient.BlockByNumber(ctx, num)
	assert.NoError(t, err)
	assert.Equal(t, block, res)
}

func TestEthClient_BlockNumber(t *testing.T) {
	ethClient, client, ctx := initClient(t)
	num := big.NewInt(1)

	// verify retry
	client.EXPECT().BlockNumber(gomock.Any()).Return(uint64(0), testErr).Times(1)
	client.EXPECT().BlockNumber(gomock.Any()).Return(num.Uint64(), nil).Times(1)

	res, err := ethClient.BlockNumber(ctx)
	assert.NoError(t, err)
	assert.Equal(t, num.Uint64(), res)
}

func TestEthClient_TransactionReceipt(t *testing.T) {
	ethClient, client, ctx := initClient(t)

	txHash := common.HexToHash(testTxHash)

	// verify retry
	client.EXPECT().TransactionReceipt(gomock.Any(), txHash).Return(nil, testErr).Times(1)
	client.EXPECT().TransactionReceipt(gomock.Any(), txHash).Return(nil, nil).Times(1)

	res, err := ethClient.TransactionReceipt(ctx, txHash)
	assert.NoError(t, err)
	assert.Nil(t, res)
}
