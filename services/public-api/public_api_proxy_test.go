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
	ctx := context.WithValue(req.Context(), botIDKey, "test-id")
	ctx = context.WithValue(ctx, botOwnerKey, "test-owner")
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
	botCfg := &config.AgentConfig{Owner: "test-owner", ID: "test-id"}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	botRemoteAddr := "1.1.1.1:1111"
	req.RemoteAddr = botRemoteAddr

	proxy := PublicAPIProxy{authenticator: authenticator}
	authenticator.EXPECT().FindContainerNameFromRemoteAddr(gomock.Any(), botRemoteAddr).Return("forta-bot-1", nil)
	authenticator.EXPECT().FindAgentByContainerName("forta-bot-1").Return(botCfg, nil)
	req, err := proxy.authenticateRequest(req)
	assert.NotNil(t, req)
	assert.NoError(t, err)

	botID, botOwner, ok := getBotFromContext(req.Context())
	assert.True(t, ok)
	assert.Equal(t, botID, "test-id")
	assert.Equal(t, botOwner, "test-owner")

	// Case 2: proxying handle alert request
	botCfg = &config.AgentConfig{Owner: "test-combiner-owner", ID: "test-combiner-id"}
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

	botID, botOwner, ok = getBotFromContext(req.Context())
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

	botID, botOwner, ok = getBotFromContext(req.Context())
	assert.False(t, ok)
	assert.Empty(t, botID)
	assert.Empty(t, botOwner)
}
