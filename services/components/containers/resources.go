package containers

import (
	"context"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/services/components/metrics"
	"github.com/sirupsen/logrus"
)

var (
	defaultPollingInterval = time.Second * 10
)

// PollDockerResources gets CPU and MEM usage for all agent containers with docker stats.
func PollDockerResources(
	ctx context.Context, dockerClient clients.DockerClient,
	msgClient clients.MessageClient,
) {
	pollingTicker := time.NewTicker(defaultPollingInterval)

	for {
		select {
		case <-ctx.Done():
			logrus.Info("stopping docker resources poller")
			return
		case <-pollingTicker.C:
			logrus.Info("polling docker resources")
			containers, err := dockerClient.GetContainers(ctx)
			if err != nil {
				logrus.WithError(err).Error("error while getting docker containers")
				continue
			}

			for _, container := range containers {
				resources, err := dockerClient.ContainerStats(ctx, container.ID)
				if err != nil {
					logrus.WithError(err).Error("error while getting container stats", container.ID)
					continue
				}

				botID, ok := container.Labels[docker.LabelFortaBotID]
				if !ok {
					continue
				}

				var (
					bytesSent uint64
					bytesRecv uint64
				)

				for _, network := range resources.NetworkStats {
					bytesSent += network.TxBytes
					bytesRecv += network.RxBytes
				}

				metrics.SendAgentMetrics(msgClient, []*protocol.AgentMetric{
					metrics.CreateResourcesMetric(
						botID, domain.MetricDockerResourcesCPU, float64(resources.CPUStats.CPUUsage.TotalUsage)),
					metrics.CreateResourcesMetric(
						botID, domain.MetricDockerResourcesMemory, float64(resources.MemoryStats.Usage)),
					metrics.CreateResourcesMetric(
						botID, domain.MetricDockerResourcesNetworkSent, float64(bytesSent)),
					metrics.CreateResourcesMetric(
						botID, domain.MetricDockerResourcesNetworkReceive, float64(bytesRecv)),
				})
			}
		}
	}
}
