package feeds

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"forta-network/forta-node/domain"
	clients "forta-network/forta-node/ethereum/mocks"
	"forta-network/forta-node/testutils"
	"forta-network/forta-node/utils"
)

func getTestTransactionFeed(t *testing.T, blockFeed BlockFeed) (*transactionFeed, *clients.MockClient) {
	blocks := make(chan *domain.BlockEvent, 1)
	txs := make(chan *domain.TransactionEvent, 1)
	ctrl := gomock.NewController(t)
	client := clients.NewMockClient(ctrl)
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
	bf := NewMockBlockFeed([]*domain.BlockEvent{
		{
			EventType: domain.EventTypeBlock,
			Block:     testutils.TestBlock(1, 2, 3),
		},
		{
			EventType: domain.EventTypeBlock,
			Block:     testutils.TestBlock(4, 5, 6, 6), // with duplicate
		},
		{
			EventType: domain.EventTypeBlock,
			Block:     testutils.TestBlock(), // empty
		},
		{
			EventType: domain.EventTypeBlock,
			Block:     testutils.TestBlock(7, 8, 9),
		},
	})

	txFeed, client := getTestTransactionFeed(t, bf)

	client.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(nil, nil).Times(9)

	var evts []*domain.TransactionEvent
	err := txFeed.ForEachTransaction(func(evt *domain.BlockEvent) error { return nil }, func(evt *domain.TransactionEvent) error {
		evts = append(evts, evt)
		return nil
	})

	assert.Equal(t, endOfBlocks, err)
	assert.Len(t, evts, 9)
}

func TestTransactionFeed_ToMessage(t *testing.T) {
	bf := NewMockBlockFeed([]*domain.BlockEvent{
		{
			EventType: domain.EventTypeBlock,
			Block:     testutils.TestBlock(1),
		},
	})

	txFeed, client := getTestTransactionFeed(t, bf)

	client.EXPECT().TransactionReceipt(gomock.Any(), gomock.Any()).Return(nil, nil)

	var result *domain.TransactionEvent
	err := txFeed.ForEachTransaction(func(evt *domain.BlockEvent) error { return nil }, func(evt *domain.TransactionEvent) error {
		result = evt
		return nil
	})
	assert.Equal(t, endOfBlocks, err)

	msg, err := result.ToMessage()
	assert.NoError(t, err)
	assert.Equal(t, result.Transaction.Hash, msg.Transaction.Hash)
}
