package json_rpc_cache

import (
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := cache{
		chains:      make(map[uint64]*chainCache),
		cacheExpire: time.Millisecond * 500,
	}

	cache.Append(events)

	blockNumber, ok := cache.Get(1, "eth_blockNumber", "[]")
	assert.True(t, ok)
	assert.Equal(t, "1", blockNumber)

	blockNumber, ok = cache.Get(2, "eth_blockNumber", "[]")
	assert.True(t, ok)
	assert.Equal(t, "101", blockNumber)

	time.Sleep(time.Second)

	blockNumber, ok = cache.Get(1, "eth_blockNumber", "[]")
	assert.False(t, ok)
	assert.Empty(t, blockNumber)
}

var events = &protocol.CombinedBlockEvents{
	Events: []*protocol.CombinedBlockEvent{
		{
			ChainID: 1,
			Block: &protocol.CombinedBlock{
				Hash:   "0xaaaa",
				Number: "1",
				Transactions: []*protocol.Transaction{
					{
						Hash: "0xbbbb",
						From: "0xcccc",
					},
				},
				Uncles: []string{"0xdddd"},
			},
			Logs: []*protocol.LogEntry{
				{
					Address: "0xcccc",
					Topics:  []string{"0xeeee"},
				},
			},
			Traces: []*protocol.Trace{
				{
					Action: &protocol.TraceAction{
						From: "0xcccc",
					},
					Result: &protocol.TraceResult{
						Address: "0xcccc",
					},
					TraceAddress: []int64{1},
				},
			},
		},
		{
			ChainID: 2,
			Block: &protocol.CombinedBlock{
				Hash:   "0xffff",
				Number: "100",
				Transactions: []*protocol.Transaction{
					{
						Hash: "0x1111",
						From: "0x2222",
					},
				},
				Uncles: []string{"0x3333"},
			},
			Logs: []*protocol.LogEntry{},
			Traces: []*protocol.Trace{
				{
					TraceAddress: []int64{2},
				},
			},
		},
		{
			ChainID: 2,
			Block: &protocol.CombinedBlock{
				Hash:   "0xfffd",
				Number: "101",
				Transactions: []*protocol.Transaction{
					{
						Hash: "0x1112",
						From: "0x2223",
					},
				},
				Uncles: []string{"0x3333"},
			},
			Logs: []*protocol.LogEntry{},
			Traces: []*protocol.Trace{
				{
					TraceAddress: []int64{1},
				},
			},
		},
	},
}
