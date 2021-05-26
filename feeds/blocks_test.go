package feeds

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"OpenZeppelin/fortify-node/domain"
	mocks "OpenZeppelin/fortify-node/ethereum/mocks"
	"OpenZeppelin/fortify-node/utils"
)

var testErr = errors.New("test")
var startHash = "0x4fc0862e76691f5312964883954d5c2db35e2b8f7a4f191775a4f50c69804a8d"
var reorgHash = "0xb9b293da464be42bbb87695c372678ea93a2ef87dc54213bbaa93bd6d8880c17"

var endOfBlocks = errors.New("end of blocks")

// mockBlockFeed is a mock block feed for tests
type mockBlockFeed struct {
	blocks []*domain.BlockEvent
}

// ForEachBlock is a test method that iterates over mocked blocks
func (bf *mockBlockFeed) ForEachBlock(handler func(evt *domain.BlockEvent) error) error {
	for _, b := range bf.blocks {
		if err := handler(b); err != nil {
			return err
		}
	}
	return endOfBlocks
}

// NewMockBlockFeed returns a new mockBlockFeed for tests
func NewMockBlockFeed(blocks []*domain.BlockEvent) *mockBlockFeed {
	return &mockBlockFeed{blocks}
}

func getTestBlockFeed(t *testing.T) (*blockFeed, *mocks.MockClient, context.Context, context.CancelFunc) {
	ctrl := gomock.NewController(t)
	client := mocks.NewMockClient(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	cache := utils.NewCache(10000)
	return &blockFeed{
		start:  big.NewInt(1),
		ctx:    ctx,
		client: client,
		cache:  cache,
	}, client, ctx, cancel
}

func blockWithParent(hash string, num int) *domain.Block {
	return &domain.Block{
		Hash:       fmt.Sprintf("0x%s%d", hash, num),
		ParentHash: hash,
		Number:     utils.BigIntToHex(big.NewInt(int64(num))),
	}
}

func blockEvent(blk *domain.Block) *domain.BlockEvent {
	return &domain.BlockEvent{
		EventType: domain.EventTypeBlock,
		Block:     blk,
	}
}

func reorgEvent(blk *domain.Block) *domain.BlockEvent {
	return &domain.BlockEvent{
		EventType: domain.EventTypeReorg,
		Block:     blk,
	}
}

func assertEvts(t *testing.T, actual []*domain.BlockEvent, expected ...*domain.BlockEvent) {
	assert.Equal(t, len(actual), len(expected), "expect same length")
	for i, exp := range expected {
		assert.Equal(t, exp, actual[i])
	}
}

func TestBlockFeed_ForEachBlock(t *testing.T) {
	bf, client, ctx, _ := getTestBlockFeed(t)

	block1 := blockWithParent(startHash, 1)
	block2 := blockWithParent(block1.Hash, 2)
	block3 := blockWithParent(block2.Hash, 3)

	//TODO: actually test that the trace part matters (this returns nil for now)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(1)).Return(block1, nil).Times(1)
	client.EXPECT().TraceBlock(ctx, utils.HexToBigInt(block1.Number)).Return(nil, nil).Times(1)

	client.EXPECT().BlockByNumber(ctx, big.NewInt(2)).Return(block2, nil).Times(1)
	client.EXPECT().TraceBlock(ctx, utils.HexToBigInt(block2.Number)).Return(nil, nil).Times(1)

	client.EXPECT().BlockByNumber(ctx, big.NewInt(3)).Return(block3, nil).Times(1)
	client.EXPECT().TraceBlock(ctx, utils.HexToBigInt(block3.Number)).Return(nil, nil).Times(1)

	count := 0
	var evts []*domain.BlockEvent
	res := bf.ForEachBlock(func(evt *domain.BlockEvent) error {
		count++
		evts = append(evts, evt)
		if count == 3 {
			return testErr
		}
		return nil
	})
	assert.Error(t, testErr, res)
	assert.Equal(t, 3, len(evts))
	assertEvts(t, evts, blockEvent(block1), blockEvent(block2), blockEvent(block3))
}

func TestBlockFeed_ForEachBlock_Cancelled(t *testing.T) {
	bf, client, ctx, cancel := getTestBlockFeed(t)

	hash1 := "0x4fc0862e76691f5312964883954d5c2db35e2b8f7a4f191775a4f50c69804a8d"
	block1 := blockWithParent(hash1, 1)

	client.EXPECT().BlockByNumber(ctx, big.NewInt(1)).Return(block1, nil).Times(1)
	client.EXPECT().TraceBlock(ctx, utils.HexToBigInt(block1.Number)).Return(nil, nil).Times(1)

	count := 0
	var evts []*domain.BlockEvent
	res := bf.ForEachBlock(func(evt *domain.BlockEvent) error {
		count++
		evts = append(evts, evt)
		cancel()
		return nil
	})
	assert.Error(t, context.Canceled, res)
	assert.Equal(t, 1, len(evts))
	assertEvts(t, evts, blockEvent(block1))
}

func TestBlockFeed_ForEachBlock_Reorg(t *testing.T) {
	bf, client, ctx, _ := getTestBlockFeed(t)

	// START
	block1 := blockWithParent(startHash, 1)

	// Different Parent!
	reorg := blockWithParent(reorgHash, 1)
	// Reorg...Its Parent is START, found common ancestry (exists in cache)
	block2 := blockWithParent(block1.Hash, 2)
	// And Continue
	block3 := blockWithParent(block2.Hash, 3)

	client.EXPECT().TraceBlock(ctx, big.NewInt(1)).Return(nil, nil).Times(1)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(1)).Return(block1, nil).Times(1)

	client.EXPECT().TraceBlock(ctx, big.NewInt(2)).Return(nil, nil).Times(1)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(2)).Return(reorg, nil).Times(1)

	client.EXPECT().BlockByHash(ctx, reorg.ParentHash).Return(block2, nil).Times(1)
	client.EXPECT().TraceBlock(ctx, big.NewInt(2)).Return(nil, nil).Times(1)

	client.EXPECT().TraceBlock(ctx, big.NewInt(3)).Return(nil, nil).Times(1)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(3)).Return(block3, nil).Times(1)
	client.EXPECT().TraceBlock(ctx, utils.HexToBigInt(block3.Number)).Return(nil, nil).Times(1)

	count := 0
	var evts []*domain.BlockEvent
	res := bf.ForEachBlock(func(evt *domain.BlockEvent) error {
		count++
		evts = append(evts, evt)
		if count == 4 {
			return testErr
		}
		return nil
	})
	assert.Error(t, testErr, res)
	assert.Equal(t, 4, count)
	assertEvts(t, evts, blockEvent(block1), blockEvent(reorg), reorgEvent(block2), blockEvent(block3))
}
