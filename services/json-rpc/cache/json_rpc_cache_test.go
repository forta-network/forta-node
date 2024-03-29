package json_rpc_cache

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	mock_registry "github.com/forta-network/forta-node/services/components/registry/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testBotCfg = &config.AgentConfig{Owner: "test-owner", ID: "test-id"}
)

func TestJsonRpcCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	blocksDataClient := mock_clients.NewMockBlocksDataClient(ctrl)
	authenticator := mock_clients.NewMockIPAuthenticator(ctrl)
	botRegistry := mock_registry.NewMockBotRegistry(ctrl)
	msgClient := mock_clients.NewMockMessageClient(ctrl)

	botRegistry.EXPECT().LoadAssignedBots().Return([]config.AgentConfig{*testBotCfg}, nil).AnyTimes()
	msgClient.EXPECT().PublishProto(gomock.Any(), gomock.Any()).AnyTimes()

	count := 0
	appended := make(chan struct{})
	blocksDataClient.EXPECT().GetBlocksData(gomock.Any()).Return(blocks, nil).Do(func(any) {
		count++
		if count == 2 {
			close(appended)
		}
	}).AnyTimes()

	botRemoteAddr := "1.1.1.1:1111"

	authenticator.EXPECT().FindAgentFromRemoteAddr(botRemoteAddr).Return(testBotCfg, nil)

	jrpCache := JsonRpcCache{
		ctx:              context.TODO(),
		botAuthenticator: authenticator,
		blocksDataClient: blocksDataClient,
		cfg: config.JsonRpcCacheConfig{
			CacheExpirePeriodSeconds: 300,
		},
		cache:       NewCache(300 * time.Second),
		botRegistry: botRegistry,
		msgClient:   msgClient,
	}

	go jrpCache.pollBlocksData()

	<-appended

	jrpReq := jsonRpcReq{
		ID:     json.RawMessage("1"),
		Method: "eth_blockNumber",
		Params: json.RawMessage("[]"),
	}
	b, err := json.Marshal(jrpReq)
	require.NoError(t, err)

	r := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	r.RemoteAddr = botRemoteAddr
	r.Header.Set("X-Forta-Chain-ID", "1")
	rw := httptest.NewRecorder()

	jrpCache.Handler().ServeHTTP(rw, r)

	require.Equal(t, 200, rw.Code)

	b = rw.Body.Bytes()
	var resp jsonRpcResp
	require.NoError(t, json.Unmarshal(b, &resp))
	require.Nil(t, resp.Error)

	assert.Equal(t, jrpReq.ID, resp.ID)
	assert.Equal(t, json.RawMessage(`"1"`), resp.Result)
}

func TestJsonRpcCache_NoAgentsAssigned(t *testing.T) {
	ctrl := gomock.NewController(t)
	blocksDataClient := mock_clients.NewMockBlocksDataClient(ctrl)
	authenticator := mock_clients.NewMockIPAuthenticator(ctrl)
	botRegistry := mock_registry.NewMockBotRegistry(ctrl)
	msgClient := mock_clients.NewMockMessageClient(ctrl)

	msgClient.EXPECT().PublishProto(gomock.Any(), gomock.Any()).AnyTimes()

	loaded := make(chan struct{})
	botRegistry.EXPECT().LoadAssignedBots().Return([]config.AgentConfig{}, nil).AnyTimes().Do(func() {
		close(loaded)
	})
	blocksDataClient.EXPECT().GetBlocksData(gomock.Any()).Return(blocks, nil).Times(0)

	jrpCache := JsonRpcCache{
		ctx:              context.TODO(),
		botAuthenticator: authenticator,
		blocksDataClient: blocksDataClient,
		cfg: config.JsonRpcCacheConfig{
			CacheExpirePeriodSeconds: 300,
		},
		cache:       NewCache(300 * time.Second),
		botRegistry: botRegistry,
		msgClient:   msgClient,
	}

	go jrpCache.pollBlocksData()
	<-loaded

	require.Equal(t, 0, jrpCache.cache.cache.ItemCount())

}
