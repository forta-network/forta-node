package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// nodeReleaseInfo defines the command output format.
type nodeReleaseInfo struct {
	CLI        *release.ReleaseSummary `json:"cli"`
	Containers *release.ReleaseSummary `json:"containers,omitempty"`
}

func handleFortaVersion(cmd *cobra.Command, args []string) error {
	dockerClient, err := clients.NewDockerClient("")
	if err != nil {
		return fmt.Errorf("failed to create the docker client: %v", err)
	}

	output, err := makeFortaVersionOutput(dockerClient)
	if err != nil {
		return nil
	}
	fmt.Println(output)
	return nil
}

func makeFortaVersionOutput(dockerClient clients.DockerClient) (string, error) {
	var info nodeReleaseInfo

	releaseSummary, ok := config.GetBuildReleaseSummary()
	if !ok {
		releaseSummary = &release.ReleaseSummary{
			Version: "custom",
		}
	}

	info.CLI = releaseSummary
	info.Containers = getReleaseInfoFromScannerContainer(dockerClient)

	b, _ := json.MarshalIndent(info, "", "  ")
	return string(b), nil
}

func getReleaseInfoFromScannerContainer(dockerClient clients.DockerClient) *release.ReleaseSummary {
	container, err := dockerClient.GetContainerByName(context.Background(), config.DockerScannerContainerName)
	if err != nil {
		return nil
	}
	containerConfig, err := dockerClient.InspectContainer(context.Background(), container.ID)
	if err != nil {
		return nil
	}
	for _, v := range containerConfig.Config.Env {
		parts := strings.Split(v, "=")
		if parts[0] != config.EnvReleaseInfo {
			continue
		}
		log.SetLevel(log.ErrorLevel)
		releaseInfo := release.ReleaseInfoFromString(parts[1])
		return release.MakeSummaryFromReleaseInfo(releaseInfo)
	}
	return nil
}
