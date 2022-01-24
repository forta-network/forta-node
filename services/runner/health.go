package runner

import (
	"fmt"
	"strconv"

	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/config"
)

func (runner *Runner) checkHealth() (reports health.Reports) {
	containers, err := runner.globalClient.GetFortaServiceContainers(runner.ctx)
	if err != nil {
		return health.Reports{
			{
				Name:    "docker",
				Status:  health.StatusDown,
				Details: fmt.Sprintf("failed to get running services: %v", err),
			},
		}
	}

	for _, container := range containers {
		name := fmt.Sprintf("container.%s", container.Names[0][1:])

		if container.State != "running" {
			reports = append(reports, &health.Report{
				Name:    name,
				Status:  health.StatusDown,
				Details: fmt.Sprintf("state: %s", container.State),
			})
			continue
		}

		reports = append(reports, &health.Report{
			Name:    name,
			Status:  health.StatusOK,
			Details: fmt.Sprintf("state: %s", container.State),
		})

		// no further checks if nats
		if container.Names[0][1:] == config.DockerNatsContainerName {
			continue
		}

		var gotReports bool
		for _, port := range container.Ports {
			if strconv.Itoa(int(port.PrivatePort)) == config.DefaultHealthPort {
				reports = append(reports, runner.healthClient.CheckHealth(name, strconv.Itoa(int(port.PublicPort)))...)
				gotReports = true
				break
			}
		}
		if gotReports {
			continue
		}
		reports = append(reports, &health.Report{
			Name:    name,
			Status:  health.StatusInfo,
			Details: "no source found",
		})
	}
	return
}
