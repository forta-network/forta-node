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
	mock_ratelimiter "github.com/forta-network/forta-node/clients/rate_limiter/mocks"
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
	ctx := context.WithValue(req.Context(), authenticatedBotKey, testBotCfg)
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

func TestPublicAPIProxy_authenticateBotRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := mock_clients.NewMockBotAuthenticator(ctrl)

	// Case 1: proxying a bot request
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	botRemoteAddr := "1.1.1.1:1111"
	req.RemoteAddr = botRemoteAddr

	proxy := PublicAPIProxy{botAuthenticator: authenticator}
	authenticator.EXPECT().FindAgentFromRemoteAddr(botRemoteAddr).Return(testBotCfg, nil)
	req, err := proxy.authenticateBotRequest(req)
	assert.NotNil(t, req)
	assert.NoError(t, err)

	bot, ok := getBotFromContext(req.Context())
	assert.True(t, ok)
	assert.Equal(t, bot.ID, "test-id")

	// Case 1: proxying an arbitrary request
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	botRemoteAddr = "1.1.1.1:1111"
	req.RemoteAddr = botRemoteAddr

	proxy = PublicAPIProxy{botAuthenticator: authenticator}
	authenticator.EXPECT().FindAgentFromRemoteAddr(botRemoteAddr).Return(nil, fmt.Errorf("can't find"))
	req, err = proxy.authenticateBotRequest(req)
	assert.NotNil(t, req)
	assert.Error(t, err)

	bot, ok = getBotFromContext(req.Context())
	assert.False(t, ok)
	assert.Nil(t, bot)
}

func TestPublicAPIProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := mock_clients.NewMockBotAuthenticator(ctrl)
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
	authenticator.EXPECT().FindAgentFromRemoteAddr(gomock.Any()).Return(nil, fmt.Errorf("can't find"))

	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)

	// case 2: authorized request
	authenticator.EXPECT().FindAgentFromRemoteAddr(gomock.Any()).Return(testBotCfg, nil)
	ratelimiter.EXPECT().ExceedsLimit(gomock.Any()).Return(false)

	resp, err = http.Get(server.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// case 2: authorized, but rate limited request
	authenticator.EXPECT().FindAgentFromRemoteAddr(gomock.Any()).Return(testBotCfg, nil)
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