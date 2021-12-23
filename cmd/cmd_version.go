package cmd

import (
	"encoding/json"

	"github.com/forta-protocol/forta-node/config"
	"github.com/spf13/cobra"
)

func handleFortaVersion(cmd *cobra.Command, args []string) error {
	releaseInfo, ok := config.GetBuildReleaseInfo()
	if !ok {
		return nil
	}
	b, _ := json.MarshalIndent(releaseInfo, "", "  ")
	cmd.Println(string(b))
	return nil
}
