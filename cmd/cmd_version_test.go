package cmd

import (
	"errors"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	testContainerID = "test-container-id"
	testReleaseInfo = "{\"fromBuild\":false,\"ipfs\":\"QmVDqk2pFVJ6joc8sNsNjsT1vzkjnYgz7X1dVw8eT1EgnR\",\"manifest\":{\"release\":{\"timestamp\":\"2023-05-12T19:11:59Z\",\"repository\":\"https://github.com/forta-network/forta-node\",\"version\":\"v0.0.1\",\"commit\":\"6c63eeee87cc99819899008ffc5db3165636bfe3\",\"services\":{\"updater\":\"disco-dev.forta.network/bafybeicvwir2cnqax74yuge2qs7qo7n7vkuqjra2wdrnhulfznxhyb6v54@sha256:2c3cee6efb77af6def70ec9e053e3e2b37300ce0f6347f5e8250f6f23a5e1467\",\"supervisor\":\"disco-dev.forta.network/bafybeicvwir2cnqax74yuge2qs7qo7n7vkuqjra2wdrnhulfznxhyb6v54@sha256:2c3cee6efb77af6def70ec9e053e3e2b37300ce0f6347f5e8250f6f23a5e1467\"},\"config\":{\"autoUpdateInHours\":24,\"deprecationPolicy\":{\"supportedVersions\":[\"v0.0.1\"],\"activatesInHours\":168}}}}}"
)

const testBothVersionsOutput = `{
  "cli": {
    "version": "custom"
  },
  "containers": {
    "commit": "6c63eeee87cc99819899008ffc5db3165636bfe3",
    "ipfs": "QmVDqk2pFVJ6joc8sNsNjsT1vzkjnYgz7X1dVw8eT1EgnR",
    "version": "v0.0.1"
  }
}`

const testCliVersionOnlyOutput = `{
  "cli": {
    "version": "custom"
  }
}`

func TestMakeFortaVersionOutput_BothVersions(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)

	dockerClient.EXPECT().GetContainerByName(gomock.Any(), config.DockerScannerContainerName).Return(&types.Container{
		ID: testContainerID,
	}, nil)
	dockerClient.EXPECT().InspectContainer(gomock.Any(), testContainerID).Return(&types.ContainerJSON{
		Config: &container.Config{
			Env: []string{fmt.Sprintf("%s=%s", config.EnvReleaseInfo, testReleaseInfo)},
		},
	}, nil)

	output, err := makeFortaVersionOutput(dockerClient)
	r.NoError(err)
	r.Equal(testBothVersionsOutput, output)
}

func TestMakeFortaVersionOutput_OnlyCLI(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)

	dockerClient.EXPECT().GetContainerByName(gomock.Any(), config.DockerScannerContainerName).Return(nil, errors.New("some error"))

	output, err := makeFortaVersionOutput(dockerClient)
	r.NoError(err)
	r.Equal(testCliVersionOnlyOutput, output)
}

func TestHandleFortaVersion(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)

	dockerClient.EXPECT().GetContainerByName(gomock.Any(), config.DockerScannerContainerName).Return(nil, errors.New("some error"))

	err := handleFortaVersionWithDockerClient(dockerClient)
	r.NoError(err)
}
