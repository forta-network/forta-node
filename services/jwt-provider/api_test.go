package jwt_provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/forta-network/forta-node/services/jwt-provider/provider"
	mock_provider "github.com/forta-network/forta-node/services/jwt-provider/provider/mocks"
)

func TestHandleJwtRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name             string
		requestBody      func() []byte
		remoteAddr       string
		mockFunc         func(mockProvider *mock_provider.MockJWTProvider)
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			name:       "successful case",
			remoteAddr: "127.0.0.1:12345",
			requestBody: func() []byte {
				b, err := json.Marshal(&CreateJWTMessage{
					Claims: map[string]interface{}{"test-claim": "claim-value"},
				})
				assert.NoError(t, err)
				return b
			},
			mockFunc: func(mockProvider *mock_provider.MockJWTProvider) {
				// Define the behavior for successful case
				mockProvider.EXPECT().CreateJWTFromIP(gomock.Any(), gomock.Any(), gomock.Any()).Return("mockJWT", nil)
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: `{"token":"mockJWT"}`,
		},
		{
			name:       "bad remote addr",
			remoteAddr: "not valid",
			requestBody: func() []byte {
				b, err := json.Marshal(&CreateJWTMessage{
					Claims: map[string]interface{}{"test-claim": "claim-value"},
				})
				assert.NoError(t, err)
				return b
			},
			mockFunc: func(mockProvider *mock_provider.MockJWTProvider) {
				// No need to define anything for provider, as we expect to fail before calling it
			},
			expectedHTTPCode: http.StatusUnauthorized,
			expectedResponse: fmt.Sprintf("can't extract ip from request: %s", "not valid"),
		},
		{
			name:       "bad create jwt message body",
			remoteAddr: "127.0.0.1:12345",
			requestBody: func() []byte {
				return []byte("bad json")
			},
			mockFunc: func(mockProvider *mock_provider.MockJWTProvider) {
				// No need to define anything for provider, as we expect to fail before calling it
			},
			expectedHTTPCode: http.StatusBadRequest,
			expectedResponse: errBadCreateMessage,
		},
		{
			name:       "failure on CreateJWTFromIP (can't find bot id from request source)",
			remoteAddr: "127.0.0.1:12345",
			requestBody: func() []byte {
				b, err := json.Marshal(&CreateJWTMessage{
					Claims: map[string]interface{}{"test-claim": "claim-value"},
				})
				assert.NoError(t, err)
				return b
			},
			mockFunc: func(mockProvider *mock_provider.MockJWTProvider) {
				// Define the behavior to make CreateJWTFromIP return ErrCannotFindBotForIP
				mockProvider.EXPECT().CreateJWTFromIP(gomock.Any(), gomock.Any(), gomock.Any()).Return("", provider.ErrCannotFindBotForIP)
			},
			expectedHTTPCode: http.StatusForbidden,
			expectedResponse: "can't find bot id from request source 127.0.0.1, err: cannot find bot for ip",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProvider := mock_provider.NewMockJWTProvider(ctrl)
			tc.mockFunc(mockProvider)

			// Create JWTAPI
			api := &JWTAPI{
				provider: mockProvider,
			}

			// Create request body
			req := httptest.NewRequest("POST", "http://localhost/create", bytes.NewBuffer(tc.requestBody()))
			req.RemoteAddr = tc.remoteAddr
			w := httptest.NewRecorder()

			api.handleJwtRequest(w, req)

			resp := w.Result()
			require.Equal(t, tc.expectedHTTPCode, resp.StatusCode)

			// Read the response
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(resp.Body)
			newStr := buf.String()

			require.Equal(t, tc.expectedResponse, newStr)
		})
	}
}
