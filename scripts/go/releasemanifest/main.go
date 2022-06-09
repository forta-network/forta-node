package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/forta-network/forta-core-go/registry"
	"github.com/forta-network/forta-core-go/release"
	"github.com/spf13/cobra"
)

const (
	defaultIpfsGateway = "https://ipfs.forta.network"
)

var (
	flagVersion      string
	flagCommitSha    string
	flagNodeImage    string
	flagIsPrerelease bool
	flagToStdout     bool
)

// version, commit sha, node image, is prerelease
func main() {
	mainCmd := &cobra.Command{
		Use:   "cmd",
		Short: "Create release manifest file",
		RunE:  handler,
	}

	mainCmd.Flags().StringVar(&flagVersion, "version", "", "")
	mainCmd.MarkFlagRequired("version")

	mainCmd.Flags().StringVar(&flagCommitSha, "commit-sha", "", "")
	mainCmd.MarkFlagRequired("commit-sha")

	mainCmd.Flags().StringVar(&flagNodeImage, "node-image", "", "")
	mainCmd.MarkFlagRequired("node-image")

	mainCmd.Flags().BoolVar(&flagIsPrerelease, "is-prerelease", false, "")

	mainCmd.Flags().BoolVar(&flagToStdout, "to-stdout", false, "")

	mainCmd.Execute()
}

func handler(cmd *cobra.Command, args []string) error {
	var newReleaseManifest release.ReleaseManifest
	if flagIsPrerelease {
		regClient, err := registry.NewDefaultClient(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to create the registry client: %v", err)
		}
		cid, err := regClient.GetScannerNodeVersion()
		if err != nil {
			return fmt.Errorf("failed to get the scanner node version: %v", err)
		}

		client, err := release.NewClient(defaultIpfsGateway)
		if err != nil {
			return fmt.Errorf("failed to create the release client: %v", err)
		}
		prevReleaseManifest, err := client.GetReleaseManifest(cmd.Context(), cid)
		if err != nil {
			return fmt.Errorf("failed to get the previous release manifest: %v", err)
		}
		newReleaseManifest.Release = prevReleaseManifest.Release
	}

	newRelease := release.Release{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Repository: "https://github.com/forta-network/forta-node",
		Version:    flagVersion,
		Commit:     flagCommitSha,
		Services: release.ReleaseServices{
			Updater:    flagNodeImage,
			Supervisor: flagNodeImage,
		},
	}

	if flagIsPrerelease {
		newReleaseManifest.Prerelease = &newRelease
	} else {
		newReleaseManifest.Release = newRelease
	}

	b, err := json.MarshalIndent(newReleaseManifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal new manifest: %v", err)
	}
	if flagToStdout {
		fmt.Println(string(b))
		return nil
	}

	if err := ioutil.WriteFile("manifest.json", b, 0755); err != nil {
		return fmt.Errorf("failed to write the new manifest file: %v", err)
	}

	return nil
}
