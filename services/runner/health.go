package runner

import (
	"fmt"
	"strconv"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-node/config"
)

// CheckServiceHealth checks the health of all node services and returns all health reports.
func (runner *Runner) CheckServiceHealth() (allReports health.Reports) {
	containers, err := runner.globalClient.GetFortaServiceContainers(runner.ctx)
	if err != nil {
		return health.Reports{
			{
				Name:    "docker",
				Status:  health.StatusDown,
				Details: err.Error(),
			},
		}
	}

	allReports = append(allReports, &health.Report{
		Name:    "forta.version",
		Status:  health.StatusInfo,
		Details: config.GetBuildReleaseInfo().Manifest.Release.Version,
	})

	for _, container := range containers {
		name := fmt.Sprintf("forta.container.%s", container.Names[0][1:])

		if container.State != "running" {
			allReports = append(allReports, &health.Report{
				Name:    name,
				Status:  health.StatusDown,
				Details: container.State,
			})
			continue
		}

		allReports = append(allReports, &health.Report{
			Name:    name,
			Status:  health.StatusOK,
			Details: container.State,
		})

		// no further checks if nats
		if container.Names[0][1:] == config.DockerNatsContainerName {
			continue
		}

		var gotReports bool
		for _, port := range container.Ports {
			if strconv.Itoa(int(port.PrivatePort)) == config.DefaultHealthPort {
				reports := runner.healthClient.CheckHealth(name, strconv.Itoa(int(port.PublicPort)))
				for _, report := range reports {
					report.Name = fmt.Sprintf("%s.%s", name, report.Name)
				}
				reports.ObfuscateDetails()
				allReports = append(allReports, reports...)
				gotReports = true
				break
			}
		}
		if gotReports {
			continue
		}
		allReports = append(allReports, &health.Report{
			Name:    name,
			Status:  health.StatusInfo,
			Details: "no source found",
		})
	}
	return
}
