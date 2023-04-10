package public_api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/security"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	mock_ratelimiter "github.com/forta-network/forta-node/clients/ratelimiter/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	testBotCfg = &config.AgentConfig{Owner: "test-owner", ID: "test-id"}
)

func _keyConstructor(t *testing.T) *keystore.Key {
	dir := t.TempDir()
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := ks.NewAccount("Forta123")
	if err != nil {
		t.Fatal(err)
	}

	key, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	if err != nil {
		t.Fatal(err)
	}

	return key
}

func TestPublicAPIProxy_setAuthBearer(t *testing.T) {
	key := _keyConstructor(t)

	// Case 1: proxying a bot request
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := context.WithValue(req.Context(), botIDKey, "test-id")
	ctx = context.WithValue(ctx, botOwnerKey, "test-owner")
	ctx = context.WithValue(ctx, isScannerKey, true)
	req = req.WithContext(ctx)

	proxy := PublicAPIProxy{Key: key}
	proxy.setAuthBearer(req)
	// parse and authenticate token
	h := req.Header.Get("Authorization")
	s := strings.Split(h, "Bearer ")
	token := s[1]

	jwtToken, err := security.VerifyScannerJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, jwtToken.Token.Claims.(jwt.MapClaims)["bot-owner"], "test-owner")
	assert.Equal(t, jwtToken.Token.Claims.(jwt.MapClaims)["bot-id"], "test-id")
}

func TestPublicAPIProxy_authenticateRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := mock_clients.NewMockIPAuthenticator(ctrl)

	// Case 1: proxying a bot request
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	botRemoteAddr := "1.1.1.1:1111"
	req.RemoteAddr = botRemoteAddr

	proxy := PublicAPIProxy{authenticator: authenticator}
	authenticator.EXPECT().FindContainerNameFromRemoteAddr(gomock.Any(), botRemoteAddr).Return("forta-bot-1", nil)
	authenticator.EXPECT().FindAgentByContainerName("forta-bot-1").Return(testBotCfg, nil)
	req, err := proxy.authenticateRequest(req)
	assert.NotNil(t, req)
	assert.NoError(t, err)

	botID, botOwner, _, ok := getBotFromContext(req.Context())
	assert.True(t, ok)
	assert.Equal(t, botID, "test-id")
	assert.Equal(t, botOwner, "test-owner")

	// Case 2: proxying handle alert request
	botCfg := &config.AgentConfig{Owner: "test-combiner-owner", ID: "test-combiner-id"}
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	remoteAddr := "1.1.1.1:1111"
	req.RemoteAddr = remoteAddr
	req.Header.Set("bot-id", botCfg.ID)
	req.Header.Set("bot-owner", botCfg.Owner)

	proxy = PublicAPIProxy{authenticator: authenticator}
	authenticator.EXPECT().FindContainerNameFromRemoteAddr(gomock.Any(), remoteAddr).Return("forta-scanner", nil)
	req, err = proxy.authenticateRequest(req)
	assert.NotNil(t, req)
	assert.NoError(t, err)

	botID, botOwner, _, ok = getBotFromContext(req.Context())
	assert.True(t, ok)
	assert.Equal(t, botID, botCfg.ID)
	assert.Equal(t, botOwner, botCfg.Owner)

	// Case 3: proxying an arbitrary request
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	botRemoteAddr = "1.1.1.1:1111"
	req.RemoteAddr = botRemoteAddr

	proxy = PublicAPIProxy{authenticator: authenticator}
	authenticator.EXPECT().FindContainerNameFromRemoteAddr(gomock.Any(), botRemoteAddr).Return("", fmt.Errorf("can't find"))
	req, err = proxy.authenticateRequest(req)
	assert.NotNil(t, req)
	assert.Error(t, err)

	botID, botOwner, _, ok = getBotFromContext(req.Context())
	assert.False(t, ok)
	assert.Empty(t, botID)
	assert.Empty(t, botOwner)
}

func TestPublicAPIProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := mock_clients.NewMockIPAuthenticator(ctrl)
	messageClient := mock_clients.NewMockMessageClient(ctrl)
	ratelimiter := mock_ratelimiter.NewMockRateLimiter(ctrl)
	messageClient.EXPECT().PublishProto(gomock.Any(), gomock.Any()).AnyTimes()
	p, _ := newPublicAPIProxy(
		context.Background(), config.PublicAPIProxyConfig{
			Url:     "https://api.forta.network",
			Headers: map[string]string{"test-header": "test-header-value"},
		}, authenticator, ratelimiter, _keyConstructor(t), messageClient,
	)

	server := httptest.NewServer(p.createPublicAPIProxyHandler())

	// case 1: unauthorized request
	authenticator.EXPECT().FindContainerNameFromRemoteAddr(gomock.Any(), gomock.Any()).Return("", fmt.Errorf("can't find"))

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)

	// case 2: authorized request
	authenticator.EXPECT().FindContainerNameFromRemoteAddr(gomock.Any(), gomock.Any()).Return("forta-bot-1", nil)
	authenticator.EXPECT().FindAgentByContainerName(gomock.Any()).Return(testBotCfg, nil)
	ratelimiter.EXPECT().ExceedsLimit(gomock.Any()).Return(false)

	resp, err = http.Get(server.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// case 2: authorized, but rate limited request
	authenticator.EXPECT().FindContainerNameFromRemoteAddr(gomock.Any(), gomock.Any()).Return("forta-bot-1", nil)
	authenticator.EXPECT().FindAgentByContainerName(gomock.Any()).Return(testBotCfg, nil)
	ratelimiter.EXPECT().ExceedsLimit(gomock.Any()).Return(true)
	resp, err = http.Get(server.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusTooManyRequests)
}

func TestPublicAPIProxy_newReverseProxy(t *testing.T) {
	// can detect bad url
	cfg := config.PublicAPIProxyConfig{Url: "xxx"}
	p := PublicAPIProxy{cfg: cfg}
	h := p.newReverseProxy()
	assert.NotNil(t, h)
}