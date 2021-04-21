package feeds

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	clients "OpenZeppelin/safe-node/clients/mocks"
	"OpenZeppelin/safe-node/testutils"
	"OpenZeppelin/safe-node/utils"
)

func getTestTransactionFeed(t *testing.T, blockFeed BlockFeed) (*transactionFeed, *clients.MockEthClient) {
	blocks := make(chan *BlockEvent, 1)
	txs := make(chan *TransactionEvent, 1)
	ctrl := gomock.NewController(t)
	client := clients.NewMockEthClient(ctrl)
	cache := utils.NewCache(10000)
	return &transactionFeed{
		ctx:       context.Background(),
		blockFeed: blockFeed,
		cache:     cache,
		txCh:      txs,
		blockCh:   blocks,
		client:    client,
		workers:   1,
	}, client
}

func TestTransactionFeed_ForEachTransaction(t *testing.T) {
	bf := NewMockBlockFeed([]*BlockEvent{
		{
			EventType: EventTypeBlock,
			Block:     testutils.TestBlock(1, 2, 3),
		},
		{
			EventType: EventTypeBlock,
			Block:     testutils.TestBlock(4, 5, 6, 6), // with duplicate
		},
		{
			EventType: EventTypeBlock,
			Block:     testutils.TestBlock(), // empty
		},
		{
			EventType: EventTypeBlock,
			Block:     testutils.TestBlock(7, 8, 9),
		},
	})

	txFeed, client := getTestTransactionFeed(t, bf)

	client.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	var evts []*TransactionEvent
	err := txFeed.ForEachTransaction(func(evt *TransactionEvent) error {
		evts = append(evts, evt)
		return nil
	})

	assert.Equal(t, endOfBlocks, err)
	assert.Len(t, evts, 9)
}
