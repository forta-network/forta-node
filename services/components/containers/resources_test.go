package containers

import (
	"context"
	"testing"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPollDockerResources(t *testing.T) {
	defaultPollingInterval = 100 * time.Millisecond

	ctrl := gomock.NewController(t)
	dockerClient := mock_clients.NewMockDockerClient(ctrl)
	msgClient := mock_clients.NewMockMessageClient(ctrl)

	ctx := context.Background()

	dockerContainerID := "test-container-id"
	agentID := "test-agent-id"

	dockerClient.EXPECT().GetContainers(ctx).Return(docker.ContainerList{
		{
			ID: dockerContainerID,
			Labels: map[string]string{
				docker.LabelFortaBotID: agentID,
			},
		},
	}, nil)

	dockerClient.EXPECT().ContainerStats(ctx, dockerContainerID).Return(&docker.ContainerResources{
		CPUStats: docker.CPUStats{
			CPUUsage: docker.CPUUsage{
				TotalUsage: 33,
			},
		},
		MemoryStats: docker.MemoryStats{
			Usage: 100,
		},
		NetworkStats: map[string]docker.NetworkStats{
			"eth0": {
				RxBytes: 123,
				TxBytes: 456,
			},
		},
	}, nil)

	doneCh := make(chan struct{})

	msgClient.EXPECT().PublishProto(messaging.SubjectMetricAgent, gomock.Any()).Do(func(v1, v2 interface{}) {
		metrics := v2.(*protocol.AgentMetricList)
		assert.Len(t, metrics.Metrics, 4)

		// CPU metric
		assert.Equal(t, agentID, metrics.Metrics[0].AgentId)
		assert.Equal(t, domain.MetricDockerResourcesCPU, metrics.Metrics[0].Name)
		assert.Equal(t, float64(33), metrics.Metrics[0].Value)

		// Memory metric
		assert.Equal(t, agentID, metrics.Metrics[1].AgentId)
		assert.Equal(t, domain.MetricDockerResourcesMemory, metrics.Metrics[1].Name)
		assert.Equal(t, float64(100), metrics.Metrics[1].Value)

		// Network bytes sent metric
		assert.Equal(t, agentID, metrics.Metrics[2].AgentId)
		assert.Equal(t, domain.MetricDockerResourcesNetworkSent, metrics.Metrics[2].Name)
		assert.Equal(t, float64(456), metrics.Metrics[2].Value)

		// Network bytes received metric
		assert.Equal(t, agentID, metrics.Metrics[3].AgentId)
		assert.Equal(t, domain.MetricDockerResourcesNetworkReceive, metrics.Metrics[3].Name)
		assert.Equal(t, float64(123), metrics.Metrics[3].Value)

		close(doneCh)
	})

	go PollDockerResources(ctx, dockerClient, msgClient)
	<-doneCh
}
