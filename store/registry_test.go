package store

import (
	"context"
	"testing"

	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/feeds"
	"github.com/forta-network/forta-node/config"
)

func Test_registryStore_FindShardIDForBot(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		ChainID: 1,
		Registry: config.RegistryConfig{
			JsonRpc: config.JsonRpcConfig{Url: "https://rpc.ankr.com/polygon"},
			IPFS: config.IPFSConfig{
				GatewayURL: "https://ipfs.forta.network",
			},
		},
		ENSConfig: config.ENSConfig{
			ContractAddress: "0x08f42fcc52a9C2F391bF507C4E8688D0b53e1bd7",
		},
	}

	ethClient, err := ethereum.NewStreamEthClient(ctx, "registry", cfg.Registry.JsonRpc.Url)
	if err != nil {
		t.Fatal(err)
	}

	bf, err := feeds.NewBlockFeed(ctx, ethClient, ethClient, feeds.BlockFeedConfig{})
	if err != nil {
		t.Fatal(err)
	}
	rs, err := NewRegistryStore(ctx, cfg, ethClient, bf)
	if err != nil {
		t.Fatal(err)
	}

	agentID := "0x80ed808b586aeebe9cdd4088ea4dea0a8e322909c0e4493c993e060e89c09ed1"
	scannerAddress := "0xA20F47966F5FF9A10b1Ec97832B50ef1f2B55fdE"
	_, err = rs.FindScannerShardIDForBot(agentID, scannerAddress)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_calculateShardID(t *testing.T) {
	type args struct {
		target uint
		shards uint
		idx    uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		// [0,0,0,1,1,1]
		{
			name: "50% sharded",
			args: args{target: 6, shards: 2, idx: 3},
			want: 1,
		},
		// [0,0,1,1,2,2]
		{
			name: "66% sharded",
			args: args{target: 6, shards: 3, idx: 3},
			want: 1,
		},
		// [0,1,2,3,4,5]
		{
			name: "max redundancy",
			args: args{target: 6, shards: 6, idx: 0},
			want: 0,
		},
		// 	[0,0,0,0,0,0]
		{
			name: "no redundancy",
			args: args{target: 6, shards: 1, idx: 4},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := calculateShardID(tt.args.target, tt.args.shards, tt.args.idx); got != tt.want {
					t.Errorf("calculateShardID() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
