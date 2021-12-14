package cmd

import (
	"github.com/forta-protocol/forta-node/cmd/updater"
	"github.com/spf13/cobra"
)

func handleFortaRun(cmd *cobra.Command, args []string) error {
	updater.Run(cfg)
	return nil
}
