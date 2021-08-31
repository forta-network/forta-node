package cmd

import (
	"github.com/forta-network/forta-node/cmd/runner"

	"github.com/spf13/cobra"
)

func handleFortaRun(cmd *cobra.Command, args []string) error {
	runner.Run(cfg)
	return nil
}
