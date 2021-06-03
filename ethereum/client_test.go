package ethereum

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"fortify-node/domain"
	mocks "fortify-node/ethereum/mocks"
)

const testBlockHash = "0x4fc0862e76691f5312964883954d5c2db35e2b8f7a4f191775a4f50c69804a8d"
const testTxHash = "0x9b9cc76d6b3b51976b1396a5b417b3bf3f4b39b8fe080e4a5aef39d02be882df"

var testErr = errors.New("test err")

func initClient(t *testing.T) (*streamEthClient, *mocks.MockrpcClient, context.Context) {
	minBackoff = 1 * time.Millisecond
	maxBackoff = 1 * time.Millisecond
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mocks.NewMockrpcClient(ctrl)

	return &streamEthClient{rpcClient: client}, client, ctx
}

func TestEthClient_BlockByHash(t *testing.T) {
	ethClient, client, ctx := initClient(t)
	hash := testBlockHash
	// verify retry
	client.EXPECT().CallContext(gomock.Any(), gomock.Any(), testBlockHash).Return(testErr).Times(1)
	client.EXPECT().CallContext(gomock.Any(), gomock.Any(), blocksByHash, testBlockHash).DoAndReturn(func(ctx context.Context, result interface{}, method string, args ...interface{}) error {
		b, _ := json.Marshal(domain.Block{Hash: hash})
		return json.Unmarshal(b, result)
	}).Times(1)

	res, err := ethClient.BlockByHash(ctx, testBlockHash)
	assert.NoError(t, err)
	assert.Equal(t, hash, res.Hash)
}
