package feeds

import (
	"context"
	"math/big"
	"testing"

	"github.com/forta-protocol/forta-node/domain"

	mocks "github.com/forta-protocol/forta-node/ethereum/mocks"
	"github.com/forta-protocol/forta-node/testutils"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

const testAlertBatchSignature = "AlertBatch(bytes32,address,uint256,uint256,uint256,uint256,uint256,string)"

// testAlertBatchTopic is the topic value for the AlertBatch event, which can be used for filtering
var testAlertBatchTopic = crypto.Keccak256Hash([]byte(testAlertBatchSignature)).Hex()

func TestLogFeed_ForEachLog(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mocks.NewMockClient(ctrl)
	addr := "0x38C1e080BeEb26eeA91932178E62987598230271"
	logs := testutils.TestLogs(0, 1, 2)

	blk := testutils.TestBlock()
	client.EXPECT().BlockByNumber(gomock.Any(), nil).Return(blk, nil).Times(1)
	client.EXPECT().GetLogs(gomock.Any(), gomock.Any()).Return([]types.Log{logs[0]}, nil).Times(1)

	client.EXPECT().BlockByNumber(gomock.Any(), big.NewInt(1)).Return(blk, nil).Times(1)
	client.EXPECT().GetLogs(gomock.Any(), gomock.Any()).Return([]types.Log{logs[1]}, nil).Times(1)

	client.EXPECT().BlockByNumber(gomock.Any(), big.NewInt(2)).Return(blk, nil).Times(1)
	client.EXPECT().GetLogs(gomock.Any(), gomock.Any()).Return([]types.Log{logs[2]}, nil).Times(1)

	lf, err := NewLogFeed(ctx, client, LogFeedConfig{
		Addresses: []string{addr},
		Topics:    [][]string{{testAlertBatchTopic}},
	})
	assert.NoError(t, err)

	var foundLogs []types.Log
	err = lf.ForEachLog(func(blk *domain.Block, logEntry types.Log) error {
		foundLogs = append(foundLogs, logEntry)
		// return early
		if len(foundLogs) == 3 {
			return context.Canceled
		}
		return nil
	}, func(blk *domain.Block) error {
		return nil
	})
	// ensure expected error is the one returned
	assert.ErrorIs(t, err, context.Canceled)

	assert.Equal(t, len(logs), len(foundLogs), "should find all logs")
	for idx, fl := range foundLogs {
		assert.Equal(t, logs[idx].TxIndex, fl.TxIndex)
		assert.Equal(t, logs[idx].TxHash.Hex(), fl.TxHash.Hex())
	}
}
