package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	nw "github.com/docker/docker/api/types/network"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/forta-network/forta-node/clients/docker"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
)

func expectGetContainer(dc *mock_clients.MockDockerClient, containerID, ipAddress string) {
	dc.EXPECT().GetContainers(gomock.Any()).Return(docker.ContainerList{{
		ID: containerID,
		NetworkSettings: &types.SummaryNetworkSettings{
			Networks: map[string]*nw.EndpointSettings{
				"test": {
					IPAddress: ipAddress,
				},
			},
		},
		Mounts: nil,
	}}, nil).Times(1)
}

func expectInspect(dc *mock_clients.MockDockerClient, containerID, botID, ipAddress string) {
	dc.EXPECT().InspectContainer(gomock.Any(), containerID).Return(&types.ContainerJSON{
		Config: &container.Config{
			Env: []string{
				fmt.Sprintf("%s=%s", config.EnvFortaBotID, botID),
			},
		},
		NetworkSettings: &types.NetworkSettings{
			Networks: map[string]*nw.EndpointSettings{
				"test": {
					IPAddress: ipAddress,
				},
			},
		},
	}, nil).Times(1)
}

func TestCreateJWT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ipAddress := "test-ip"
	botID := "test-bot-id"
	containerID := "container-id"
	jwtToken := "jwt"
	claims := map[string]interface{}{"test-claim": "claim-value"}
	ctx := context.Background()

	testCases := []struct {
		name          string
		mockFunc      func(dc *mock_clients.MockDockerClient)
		expectedError error
	}{
		{
			name: "successful case",
			mockFunc: func(dc *mock_clients.MockDockerClient) {
				expectGetContainer(dc, containerID, ipAddress)
				expectInspect(dc, containerID, botID, ipAddress)
			},
			expectedError: nil,
		},
		{
			name: "get containers error",
			mockFunc: func(dc *mock_clients.MockDockerClient) {
				dc.EXPECT().GetContainers(ctx).Return(nil, errors.New("test err"))
			},
			expectedError: ErrCannotFindBotForIP,
		},
		{
			name: "inspect container error",
			mockFunc: func(dc *mock_clients.MockDockerClient) {
				expectGetContainer(dc, containerID, ipAddress)
				dc.EXPECT().InspectContainer(ctx, containerID).Return(nil, errors.New("test err"))
			},
			expectedError: ErrCannotFindBotForIP,
		},
		{
			name: "non-matching inspect case",
			mockFunc: func(dc *mock_clients.MockDockerClient) {
				expectGetContainer(dc, containerID, ipAddress)
				expectInspect(dc, containerID, "other", ipAddress)
			},
			expectedError: ErrCannotFindBotForIP,
		},
		{
			name: "non-matching ip case",
			mockFunc: func(dc *mock_clients.MockDockerClient) {
				expectGetContainer(dc, containerID, "other")
			},
			expectedError: ErrCannotFindBotForIP,
		},
		{
			name: "error with jwt case",
			mockFunc: func(dc *mock_clients.MockDockerClient) {
				expectGetContainer(dc, containerID, ipAddress)
				expectInspect(dc, containerID, botID, ipAddress)
			},
			expectedError: errors.New("bad jwt creation"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDockerClient := mock_clients.NewMockDockerClient(ctrl)

			// Define your JWTProvider
			jp := &jwtProvider{
				cfg:          config.Config{},
				dockerClient: mockDockerClient,
				jwtCreatorFunc: func(key *keystore.Key, claims map[string]interface{}) (string, error) {
					// Mock the JWT creation function here
					if tc.expectedError != nil {
						return "", tc.expectedError
					}
					return jwtToken, nil
				},
			}

			tc.mockFunc(mockDockerClient)

			// Call the method
			jwt, err := jp.CreateJWT(context.Background(), ipAddress, claims)

			if tc.expectedError != nil {
				assert.ErrorIs(t, tc.expectedError, err)
			} else {
				assert.Equal(t, jwtToken, jwt)
			}
		})
	}
}
