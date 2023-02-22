package store

import (
	"context"
	"math/big"
	"sync"
	"testing"

	"github.com/forta-network/forta-core-go/manifest"
	mock_manifest2 "github.com/forta-network/forta-core-go/manifest/mocks"
	"github.com/forta-network/forta-core-go/registry"
	mock_registry "github.com/forta-network/forta-core-go/registry/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
			name: "2 shards, 3 target",
			args: args{target: 3, shards: 2, idx: 3},
			want: 1,
		},
		// [0,1]
		{
			name: "2 shards, 1 target",
			args: args{target: 1, shards: 2, idx: 1},
			want: 1,
		},
		// [0,0,1,1,2,2]
		{
			name: "3 shards, 2 target",
			args: args{target: 2, shards: 3, idx: 3},
			want: 1,
		},
		// [0,1,2,3,4,5]
		{
			name: "6 shards, 1 target",
			args: args{target: 1, shards: 6, idx: 5},
			want: 5,
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
				if got := calculateShardID(tt.args.target, tt.args.idx); got != tt.want {
					t.Errorf("calculateShardID() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_registryStore_FindScannerShardIDForBot(t *testing.T) {
	type fields struct {
		mc  func(t *testing.T) *mock_manifest2.MockClient
		rc  func(t *testing.T) *mock_registry.MockClient
		cfg config.Config
		mu  sync.Mutex
	}
	type args struct {
		agentID        string
		scannerAddress string
	}

	var (
		mockImageRef1 = "disco.forta.network/bafybeicswtrmgplg3npq5aiqafre3bqgpzaj44x2dk5a4mrnzuuf7n5zdq@sha256:7084fe732e1874feaf8f44cbf224c06ab0d7ae214c759003b287b32c9cd43100"

		// test case 1
		mockShards1  uint = 1
		mockAssigns1 uint = 2

		// test case 2
		mockShards2    uint = 2
		mockAssigns2   uint = 4
		mockScannerIdx      = big.NewInt(int64(2))
	)

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantShardID uint
		wantShards  uint
		wantTarget  uint
		wantErr     bool
	}{
		{
			name: "can calculate correctly for unsharded bot",
			args: args{agentID: "0xbot1", scannerAddress: "0xscanner1"},
			fields: fields{
				mc: func(t *testing.T) *mock_manifest2.MockClient {
					ctrl := gomock.NewController(t)
					m := mock_manifest2.NewMockClient(ctrl)

					m.EXPECT().GetAgentManifest(gomock.Any(), gomock.Any()).Return(
						&manifest.SignedAgentManifest{
							Manifest: &manifest.AgentManifest{ImageReference: &mockImageRef1},
						}, nil,
					).Times(2)
					return m
				},
				rc: func(t *testing.T) *mock_registry.MockClient {
					ctrl := gomock.NewController(t)
					rc := mock_registry.NewMockClient(ctrl)
					rc.EXPECT().GetAgent("0xbot1").Return(&registry.Agent{Manifest: "QmfFo3bw3QnuHEMPnsWmxHoHzNYd98vuEjGQJYVg9knYBm"}, nil)
					rc.EXPECT().NumScannersForByChain(gomock.Any(), gomock.Any()).Return(big.NewInt(int64(mockAssigns1)), nil)
					return rc
				},
				cfg: config.Config{},
				mu:  sync.Mutex{},
			},
			wantErr:     false,
			wantShards:  mockShards1,
			wantTarget:  mockAssigns1,
			wantShardID: 0,
		},
		{
			name: "can calculate correctly for sharded bot",
			args: args{agentID: "0xbot1", scannerAddress: "0xscanner1"},
			fields: fields{
				mc: func(t *testing.T) *mock_manifest2.MockClient {
					ctrl := gomock.NewController(t)
					m := mock_manifest2.NewMockClient(ctrl)

					m.EXPECT().GetAgentManifest(gomock.Any(), gomock.Any()).Return(
						&manifest.SignedAgentManifest{
							Manifest: &manifest.AgentManifest{
								ImageReference: &mockImageRef1, ChainSettings: map[string]manifest.AgentChainSettings{
									keyDefaultChainSetting: {
										Shards: mockShards2,
									},
								},
							},
						}, nil,
					).Times(2)
					return m
				},
				rc: func(t *testing.T) *mock_registry.MockClient {
					ctrl := gomock.NewController(t)
					rc := mock_registry.NewMockClient(ctrl)
					rc.EXPECT().GetAgent("0xbot1").Return(&registry.Agent{Manifest: "QmfFo3bw3QnuHEMPnsWmxHoHzNYd98vuEjGQJYVg9knYBm"}, nil)
					rc.EXPECT().NumScannersForByChain(gomock.Any(), gomock.Any()).Return(big.NewInt(int64(mockAssigns2)), nil)
					rc.EXPECT().IndexOfAssignedScannerByChain(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockScannerIdx, nil)
					return rc
				},
				cfg: config.Config{},
				mu:  sync.Mutex{},
			},
			wantErr:     false,
			wantShards:  mockShards2,
			wantTarget:  2,
			wantShardID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rs, _ := NewRegistryStoreFromClients(context.Background(), tt.fields.cfg, tt.fields.mc(t), tt.fields.rc(t))
				gotShardID, gotShards, gotTarget, err := rs.FindScannerShardIDForBot(tt.args.agentID, tt.args.scannerAddress)
				if (err != nil) != tt.wantErr {
					t.Errorf("FindScannerShardIDForBot() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotShardID != tt.wantShardID {
					t.Errorf("FindScannerShardIDForBot() gotShardID = %v, want %v", gotShardID, tt.wantShardID)
				}
				if gotShards != tt.wantShards {
					t.Errorf("FindScannerShardIDForBot() gotShards = %v, want %v", gotShards, tt.wantShards)
				}
				if gotTarget != tt.wantTarget {
					t.Errorf("FindScannerShardIDForBot() gotTarget = %v, want %v", gotTarget, tt.wantTarget)
				}
			},
		)
	}
}

func Test_registryStore_GetAgentsIfChanged(t *testing.T) {
	type fields struct {
		mc  func(t *testing.T) *mock_manifest2.MockClient
		rc  func(t *testing.T) *mock_registry.MockClient
		cfg config.Config
		mu  sync.Mutex
	}
	type args struct {
		scannerAddress string
	}


	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*config.AgentConfig
		want1   bool
		wantErr bool
	}{
		{
			name: "scanner with no nodes",
			args: args{scannerAddress: "0xscanner1"},
			fields: fields{
				mc: func(t *testing.T) *mock_manifest2.MockClient {
					ctrl := gomock.NewController(t)
					m := mock_manifest2.NewMockClient(ctrl)
					return m
				},
				rc: func(t *testing.T) *mock_registry.MockClient {
					ctrl := gomock.NewController(t)
					rc := mock_registry.NewMockClient(ctrl)
					rc.EXPECT().GetAssignmentHash(gomock.Any()).Return(&registry.AssignmentHash{}, nil)
					rc.EXPECT().IsEnabledScanner(gomock.Any()).Return(true, nil)
					rc.EXPECT().PegLatestBlock()
					rc.EXPECT().ResetOpts()

					rc.EXPECT().ForEachAssignedAgent(gomock.Any(),gomock.Any()).Return(nil)
					return rc
				},
				cfg: config.Config{},
				mu:  sync.Mutex{},
			},
			wantErr: false,
			want:    nil,
			want1:   true,
		},
		{
			name: "can detect disabled scanner",
			args: args{scannerAddress: "0xscanner1"},
			fields: fields{
				mc: func(t *testing.T) *mock_manifest2.MockClient {
					ctrl := gomock.NewController(t)
					m := mock_manifest2.NewMockClient(ctrl)
					return m
				},
				rc: func(t *testing.T) *mock_registry.MockClient {
					ctrl := gomock.NewController(t)
					rc := mock_registry.NewMockClient(ctrl)
					rc.EXPECT().GetAssignmentHash(gomock.Any()).Return(&registry.AssignmentHash{}, nil)
					rc.EXPECT().IsEnabledScanner(gomock.Any()).Return(false, nil)
					return rc
				},
				cfg: config.Config{},
				mu:  sync.Mutex{},
			},
			wantErr: false,
			want:    []*config.AgentConfig{},
			want1:   true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rs, _ := NewRegistryStoreFromClients(context.Background(), tt.fields.cfg, tt.fields.mc(t), tt.fields.rc(t))
				got, got1, err := rs.GetAgentsIfChanged(tt.args.scannerAddress)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAgentsIfChanged() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !assert.Equal(t, got, tt.want) {
					t.Errorf("GetAgentsIfChanged() got = %v, want %v", got, tt.want)
				}
				if got1 != tt.want1 {
					t.Errorf("GetAgentsIfChanged() got1 = %v, want %v", got1, tt.want1)
				}
			},
		)
	}
}
