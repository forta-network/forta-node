package jwt_provider

import (
	"bytes"
	"encoding/json"
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
		mockFunc         func(mockProvider *mock_provider.MockJWTProvider)
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			name: "successful case",
			requestBody: func() []byte {
				b, err := json.Marshal(&CreateJWTMessage{
					Claims: map[string]interface{}{"test-claim": "claim-value"},
				})
				assert.NoError(t, err)
				return b
			},
			mockFunc: func(mockProvider *mock_provider.MockJWTProvider) {
				// Define the behavior for successful case
				mockProvider.EXPECT().CreateJWT(gomock.Any(), gomock.Any(), gomock.Any()).Return("mockJWT", nil)
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: `{"token":"mockJWT"}`,
		},
		{
			name: "bad create jwt message body",
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
			name: "failure on CreateJWT (can't find bot id from request source)",
			requestBody: func() []byte {
				b, err := json.Marshal(&CreateJWTMessage{
					Claims: map[string]interface{}{"test-claim": "claim-value"},
				})
				assert.NoError(t, err)
				return b
			},
			mockFunc: func(mockProvider *mock_provider.MockJWTProvider) {
				// Define the behavior to make CreateJWT return ErrCannotFindBotForIP
				mockProvider.EXPECT().CreateJWT(gomock.Any(), gomock.Any(), gomock.Any()).Return("", provider.ErrCannotFindBotForIP)
			},
			expectedHTTPCode: http.StatusForbidden,
			expectedResponse: "can't find bot id from request source 127.0.0.1, err: cannot find bot for ip",
		},
		// Add more test cases as needed
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
			req.RemoteAddr = "127.0.0.1:12345" // Add remote address to the request
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
