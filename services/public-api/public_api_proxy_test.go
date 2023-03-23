package public_api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/security"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPublicAPIProxy_setAuthBearer(t *testing.T) {
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

	// Case 1: proxying a bot request
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	botCfg := &config.AgentConfig{Owner: "test-owner", ID: "test-id"}
	ctx := context.WithValue(req.Context(), authenticatedBotKey, botCfg)
	req = req.WithContext(ctx)

	proxy := PublicAPIProxy{Key: key}
	proxy.setAuthBearer(req)
	// parse and authenticate token
	h := req.Header.Get("Authorization")
	s := strings.Split(h, "Bearer ")
	token := s[1]

	jwtToken, err := security.VerifyScannerJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, jwtToken.Token.Claims.(jwt.MapClaims)["owner"], "test-owner")
	assert.Equal(t, jwtToken.Token.Claims.(jwt.MapClaims)["bot-id"], "test-id")
}

func TestPublicAPIProxy_authenticateBotRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticator := mock_clients.NewMockBotAuthenticator(ctrl)

	// Case 1: proxying a bot request
	botCfg := &config.AgentConfig{Owner: "test-owner", ID: "test-id"}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	botRemoteAddr := "1.1.1.1:1111"
	req.RemoteAddr = botRemoteAddr

	proxy := PublicAPIProxy{botAuthenticator: authenticator}
	authenticator.EXPECT().FindAgentFromRemoteAddr(botRemoteAddr).Return(botCfg, true)
	req, ok := proxy.authenticateBotRequest(req)
	assert.NotNil(t, req)
	assert.True(t, ok)

	bot, ok := getBotFromContext(req.Context())
	assert.True(t, ok)
	assert.Equal(t, bot.ID, "test-id")

	// Case 1: proxying an arbitrary request
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	botRemoteAddr = "1.1.1.1:1111"
	req.RemoteAddr = botRemoteAddr

	proxy = PublicAPIProxy{botAuthenticator: authenticator}
	authenticator.EXPECT().FindAgentFromRemoteAddr(botRemoteAddr).Return(nil, false)
	req, ok = proxy.authenticateBotRequest(req)
	assert.NotNil(t, req)
	assert.False(t, ok)

	bot, ok = getBotFromContext(req.Context())
	assert.False(t, ok)
	assert.Nil(t, bot)
}
