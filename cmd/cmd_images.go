package cmd

import (
	"github.com/forta-network/forta-node/config"

	"github.com/spf13/cobra"
)

func handleFortaImages(cmd *cobra.Command, args []string) error {
	cmd.Println("Use images:", config.UseDockerImages)
	cmd.Println("Supervisor:", config.DockerSupervisorImage)
	cmd.Println("Updater:", config.DockerUpdaterImage)
	return nil
}
