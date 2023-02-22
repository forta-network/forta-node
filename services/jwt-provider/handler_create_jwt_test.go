package jwt_provider

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/security"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func mockJWTProvider(t *testing.T) (*JWTProvider, *mock_clients.MockDockerClient) {
	ctrl := gomock.NewController(t)

	mockDockerClient := mock_clients.NewMockDockerClient(ctrl)

	dir := t.TempDir()
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := ks.NewAccount("Forta123")
	if err != nil {
		t.Fatal(err)
	}

	mockKey, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	if err != nil {
		t.Fatal(err)
	}

	return &JWTProvider{
		cfg:          &JWTProviderConfig{Key: mockKey},
		dockerClient: mockDockerClient,
	}, mockDockerClient

}

func TestJWTProvider_jwtHandler(t *testing.T) {
	j, mockDockerClient := mockJWTProvider(t)

	// test context
	ctx := context.Background()

	//
	// Test Case 1: Bot can retrieve JWT
	//
	mockBotIP := "172.0.1.2"
	mockBotRemoteAddr := mockBotIP + ":8423"
	mockBotID := "1"
	mockBotContainerInfo := []types.Container{
		{
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"local": {IPAddress: mockBotIP},
				},
			},
		},
	}
	mockBotContainerJSON := &types.ContainerJSON{
		Config: &container.Config{
			Env: []string{
				envPrefix + mockBotID,
			},
		},
	}

	mockDockerClient.EXPECT().GetContainers(gomock.Any()).Return(mockBotContainerInfo, nil)
	mockDockerClient.EXPECT().InspectContainer(gomock.Any(), gomock.Any()).Return(mockBotContainerJSON, nil)

	_, err := j.doCreateJWT(ctx, mockBotRemoteAddr, nil)
	assert.NoError(t, err)

	//
	// Test Case 2: Unknown sources can't get JWT
	//
	mockBotIP = "172.0.1.2"
	mockBotRemoteAddr = mockBotIP + ":8423"
	mockBotID = "1"
	mockBotContainerInfo = []types.Container{
		{
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{},
			},
		},
	}
	mockBotContainerJSON = &types.ContainerJSON{
		Config: &container.Config{
			Env: []string{
				envPrefix + mockBotID,
			},
		},
	}

	mockDockerClient.EXPECT().GetContainers(gomock.Any()).Return(mockBotContainerInfo, nil)

	_, err = j.doCreateJWT(ctx, mockBotRemoteAddr, nil)
	assert.Error(t, err)

	//
	// Test Case 3: Source with bad remote address
	//
	mockBotIP = "www.X.Y.Z"
	mockBotRemoteAddr = mockBotIP
	mockBotID = "1"
	mockBotContainerInfo = []types.Container{
		{
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{},
			},
		},
	}
	mockBotContainerJSON = &types.ContainerJSON{
		Config: &container.Config{
			Env: []string{
				envPrefix + mockBotID,
			},
		},
	}

	_, err = j.doCreateJWT(ctx, mockBotRemoteAddr, nil)
	assert.Error(t, err)

	//
	// Test Case 4: Source with bad remote address
	//
	mockBotIP = "172.0.1.2"
	mockBotRemoteAddr = mockBotIP + ":8423"
	mockBotID = "1"
	mockBotContainerInfo = []types.Container{
		{
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"local": {IPAddress: mockBotIP},
				},
			},
		},
	}
	mockBotContainerJSON = &types.ContainerJSON{
		Config: &container.Config{
			Env: []string{
				envPrefix + mockBotID,
			},
		},
	}

	j2, mockDockerClient2 := mockJWTProvider(t)
	j2.cfg.Key = nil
	mockDockerClient2.EXPECT().GetContainers(gomock.Any()).Return(mockBotContainerInfo, nil)
	mockDockerClient2.EXPECT().InspectContainer(gomock.Any(), gomock.Any()).Return(mockBotContainerJSON, nil)

	_, err = j2.doCreateJWT(ctx, mockBotRemoteAddr, nil)
	assert.Error(t, err)

}

func TestJWTProvider_createJWTHandler(t *testing.T) {
	j, mockDockerClient := mockJWTProvider(t)

	body := CreateJWTMessage{map[string]interface{}{"test-claim": "success"}}
	b := bytes.NewBuffer([]byte{})
	_ = json.NewEncoder(b).Encode(body)
	req := httptest.NewRequest(http.MethodPost, "/create", b)
	w := httptest.NewRecorder()

	mockBotID := "1"
	mockBotIP, _, _ := net.SplitHostPort(req.RemoteAddr)
	mockBotContainerInfo := []types.Container{
		{
			NetworkSettings: &types.SummaryNetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"local": {
						IPAddress: mockBotIP,
					},
				},
			},
		},
	}
	mockBotContainerJSON := &types.ContainerJSON{
		Config: &container.Config{
			Env: []string{
				envPrefix + mockBotID,
			},
		},
	}

	// Test case 1: can retrieve token
	mockDockerClient.EXPECT().GetContainers(gomock.Any()).Return(mockBotContainerInfo, nil)
	mockDockerClient.EXPECT().InspectContainer(gomock.Any(), gomock.Any()).Return(mockBotContainerJSON, nil)
	j.createJWTHandler(w, req)

	// Test case 2: bad request body
	b2 := bytes.NewBuffer([]byte("xxxxxx"))
	req2 := httptest.NewRequest(http.MethodPost, "/create", b2)
	j.createJWTHandler(w, req2)

	// Test case 3: can not retrieve token, bad source
	body3 := CreateJWTMessage{map[string]interface{}{"test-claim": "success"}}
	b3 := bytes.NewBuffer([]byte{})
	_ = json.NewEncoder(b3).Encode(body3)
	req3 := httptest.NewRequest(http.MethodPost, "/create", b3)
	req3.RemoteAddr = "1.1.1.1:4444"
	w3 := httptest.NewRecorder()

	mockDockerClient.EXPECT().GetContainers(gomock.Any()).Return(mockBotContainerInfo, nil)

	j.createJWTHandler(w3, req3)
}
