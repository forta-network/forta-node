package feeds

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	clients "OpenZeppelin/zephyr-node/clients/mocks"
	"OpenZeppelin/zephyr-node/utils"
)

var testErr = errors.New("test")
var startHash = "0x4fc0862e76691f5312964883954d5c2db35e2b8f7a4f191775a4f50c69804a8d"
var reorgHash = "0xb9b293da464be42bbb87695c372678ea93a2ef87dc54213bbaa93bd6d8880c17"

var endOfBlocks = errors.New("end of blocks")

// mockBlockFeed is a mock block feed for tests
type mockBlockFeed struct {
	blocks []*BlockEvent
}

// ForEachBlock is a test method that iterates over mocked blocks
func (bf *mockBlockFeed) ForEachBlock(handler func(evt *BlockEvent) error) error {
	for _, b := range bf.blocks {
		if err := handler(b); err != nil {
			return err
		}
	}
	return endOfBlocks
}

// NewMockBlockFeed returns a new mockBlockFeed for tests
func NewMockBlockFeed(blocks []*BlockEvent) *mockBlockFeed {
	return &mockBlockFeed{blocks}
}

func getTestBlockFeed(t *testing.T) (*blockFeed, *clients.MockEthClient, context.Context, context.CancelFunc) {
	ctrl := gomock.NewController(t)
	client := clients.NewMockEthClient(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	cache := utils.NewCache(10000)
	return &blockFeed{
		start:  big.NewInt(1),
		ctx:    ctx,
		client: client,
		cache:  cache,
	}, client, ctx, cancel
}

func blockWithParent(hash string, num int) *types.Block {
	return types.NewBlockWithHeader(&types.Header{
		ParentHash: common.HexToHash(hash),
		Number:     big.NewInt(int64(num)),
	})
}

func blockEvent(blk *types.Block) *BlockEvent {
	return &BlockEvent{
		EventType: EventTypeBlock,
		Block:     blk,
	}
}

func reorgEvent(blk *types.Block) *BlockEvent {
	return &BlockEvent{
		EventType: EventTypeReorg,
		Block:     blk,
	}
}

func assertEvts(t *testing.T, actual []*BlockEvent, expected ...*BlockEvent) {
	assert.Equal(t, len(actual), len(expected), "expect same length")
	for i, exp := range expected {
		assert.Equal(t, exp, actual[i])
	}
}

func TestBlockFeed_ForEachBlock(t *testing.T) {
	bf, client, ctx, _ := getTestBlockFeed(t)

	block1 := blockWithParent(startHash, 1)
	block2 := blockWithParent(block1.Hash().Hex(), 2)
	block3 := blockWithParent(block2.Hash().Hex(), 3)

	client.EXPECT().BlockByNumber(ctx, big.NewInt(1)).Return(block1, nil).Times(1)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(2)).Return(block2, nil).Times(1)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(3)).Return(block3, nil).Times(1)

	count := 0
	var evts []*BlockEvent
	res := bf.ForEachBlock(func(evt *BlockEvent) error {
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
	block1 := types.NewBlockWithHeader(&types.Header{
		ParentHash: common.HexToHash(hash1),
	})

	client.EXPECT().BlockByNumber(ctx, big.NewInt(1)).Return(block1, nil).Times(1)

	count := 0
	var evts []*BlockEvent
	res := bf.ForEachBlock(func(evt *BlockEvent) error {
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
	block2 := blockWithParent(block1.Hash().Hex(), 2)
	// And Continue
	block3 := blockWithParent(block2.Hash().Hex(), 3)

	client.EXPECT().BlockByNumber(ctx, big.NewInt(1)).Return(block1, nil).Times(1)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(2)).Return(reorg, nil).Times(1)
	client.EXPECT().BlockByHash(ctx, reorg.ParentHash()).Return(block2, nil).Times(1)
	client.EXPECT().BlockByNumber(ctx, big.NewInt(3)).Return(block3, nil).Times(1)

	count := 0
	var evts []*BlockEvent
	res := bf.ForEachBlock(func(evt *BlockEvent) error {
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
