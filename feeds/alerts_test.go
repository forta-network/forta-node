package feeds

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/forta-network/forta-node/domain"

	"github.com/forta-network/forta-node/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

type fakeLogFeed struct {
	logs []types.Log
}

func (lf *fakeLogFeed) ForEachLog(blockHandler func(blk *domain.Block) error, handler func(logEntry types.Log) error) error {
	for _, l := range lf.logs {
		if err := handler(l); err != nil {
			return err
		}
	}
	return nil
}

func TestAlertFeed_ForEachAlert(t *testing.T) {
	ctx := context.Background()
	scanner := "0x3f88c2b3e267e6b8e9dE017cdB47a59aC9Ecb284"
	ref := "QmZRibQnRBn8p7vvV1cWmg3vSnFTx6SWHDX8f4mSb2XwCC"

	// test data is from https://goerli.etherscan.io/tx/0x2a3e2ad90270fcc652abd46f132dbd5ff068c15133dd9dfe5660ad350a304582
	alertData := `
	0x0000000000000000000000000000000000000000000000000000000000c7a38d
	0000000000000000000000000000000000000000000000000000000000c7a38f
	0000000000000000000000000000000000000000000000000000000000000001
	0000000000000000000000000000000000000000000000000000000000000003
	00000000000000000000000000000000000000000000000000000000000000a0
	000000000000000000000000000000000000000000000000000000000000002e
	516d5a526962516e52426e3870377676563163576d673376536e465478365357
	4844583866346d53623258774343000000000000000000000000000000000000`

	data := strings.ReplaceAll(alertData, "\n", "")
	data = strings.ReplaceAll(data, "\t", "")
	dataBytes, err := hexutil.Decode(data)
	assert.NoError(t, err)

	log := types.Log{
		Address: common.HexToAddress("0x38C1e080BeEb26eeA91932178E62987598230271"),
		Topics: []common.Hash{
			common.HexToHash(AlertBatchTopic),
			common.HexToHash("0x0000000000000000000000003f88c2b3e267e6b8e9de017cdb47a59ac9ecb284"),
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		TxHash: common.HexToHash("0x0"),
		Data:   dataBytes,
	}
	lf := &fakeLogFeed{logs: []types.Log{log}}
	af, err := NewAlertFeed(ctx, lf)
	assert.NoError(t, err)

	var res *contracts.AlertsAlertBatch
	err = af.ForEachAlert(func(blk *domain.Block) error {
		return nil
	}, func(logEntry types.Log, batch *contracts.AlertsAlertBatch) error {
		res = batch
		return nil
	})
	assert.NoError(t, err)

	assert.NotNil(t, res)
	assert.Equal(t, scanner, res.Scanner.Hex())
	assert.Equal(t, big.NewInt(1), res.AlertCount)
	assert.Equal(t, big.NewInt(3), res.MaxSeverity)
	assert.Equal(t, big.NewInt(1), res.ChainId)
	assert.Equal(t, big.NewInt(13083533), res.BlockStart)
	assert.Equal(t, big.NewInt(13083535), res.BlockEnd)
	assert.Equal(t, ref, res.Ref)
}
