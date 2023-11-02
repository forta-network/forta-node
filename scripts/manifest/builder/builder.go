package builder

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/forta-network/forta-core-go/release"
	"gopkg.in/yaml.v3"
)

const (
	beginReleaseConfiguration = "# @begin release_config"
	endReleaseConfiguration   = "# @end release_config"
)

// BuildManifest builds a release manifest from given inputs.
func BuildManifest(version, commitSha, nodeImage, releaseNotes string) (string, error) {
	return BuildManifestWithTimestamp(time.Now(), version, commitSha, nodeImage, releaseNotes)
}

type deprecationPolicy struct {
	SupportedVersions []string `yaml:"supportedVersions"`
	ActivatesInHours  int      `yaml:"activatesInHours"`
}

// BuildManifestWithTimestamp builds a release manifest from given inputs.
func BuildManifestWithTimestamp(ts time.Time, version, commitSha, nodeImage, releaseNotes string) (string, error) {
	var releaseInfo release.Release

	releaseInfo.Timestamp = ts.UTC().Format(time.RFC3339)
	releaseInfo.Repository = "https://github.com/forta-network/forta-node"
	releaseInfo.Version = version
	releaseInfo.Commit = commitSha
	releaseInfo.Services = release.ReleaseServices{
		Updater:    nodeImage,
		Supervisor: nodeImage,
	}

	releaseConfig, err := parseReleaseConfig(releaseNotes)
	if err != nil {
		return "", err
	}

	// set config defaults

	if releaseConfig.DeprecationPolicy.ActivatesInHours == 0 {
		releaseConfig.DeprecationPolicy.ActivatesInHours = release.DefaultDeprecationHours
	}
	if len(releaseConfig.DeprecationPolicy.SupportedVersions) == 0 {
		releaseConfig.DeprecationPolicy.SupportedVersions = []string{version}
	}
	if releaseConfig.AutoUpdateInHours == 0 {
		releaseConfig.AutoUpdateInHours = release.DefaultAutoUpdateHours
	}

	releaseInfo.Config = releaseConfig

	b, err := json.MarshalIndent(&release.ReleaseManifest{
		Release: releaseInfo,
	}, "", "  ") // indentation: two spaces
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func parseReleaseConfig(releaseNotes string) (release.ReleaseConfig, error) {
	parts := strings.Split(releaseNotes, beginReleaseConfiguration)
	if len(parts) != 2 {
		return release.ReleaseConfig{}, nil
	}
	parts = strings.Split(parts[1], endReleaseConfiguration)
	if len(parts) != 2 {
		return release.ReleaseConfig{}, nil
	}
	configStr := parts[0]
	if len(configStr) == 0 {
		return release.ReleaseConfig{}, nil
	}
	var releaseConfig release.ReleaseConfig
	if err := yaml.Unmarshal([]byte(configStr), &releaseConfig); err != nil {
		return release.ReleaseConfig{}, err
	}
	return releaseConfig, nil
}
