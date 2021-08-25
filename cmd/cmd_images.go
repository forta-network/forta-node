package cmd

import (
	"forta-network/forta-node/config"

	"github.com/spf13/cobra"
)

func handleFortaImages(cmd *cobra.Command, args []string) error {
	cmd.Println("Use containers:", config.UseDockerContainers)
	cmd.Println("Scanner:", config.DockerScannerContainerImage)
	cmd.Println("Query:", config.DockerQueryContainerImage)
	cmd.Println("Proxy:", config.DockerJSONRPCProxyContainerImage)
	return nil
}
