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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testBotCfg = &config.AgentConfig{Owner: "test-owner", ID: "test-id"}
)

func TestJsonRpcCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	eventsClient := mock_clients.NewMockCombinedBlockEventsClient(ctrl)
	authenticator := mock_clients.NewMockIPAuthenticator(ctrl)

	count := 0
	appended := make(chan struct{})
	eventsClient.EXPECT().GetCombinedBlockEvents(gomock.Any()).Return(events, nil).Do(func(any) {
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
		cbeClient:        eventsClient,
		cfg: config.JsonRpcCacheConfig{
			CacheExpirePeriodSeconds: 300,
		},
		cache: NewCache(300 * time.Second),
	}

	go jrpCache.pollEvents()

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
