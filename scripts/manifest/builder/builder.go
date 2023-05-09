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

	var err error
	releaseInfo.DeprecationPolicy, err = parseDeprecationPolicy(releaseNotes)
	if err != nil {
		return "", err
	}

	// set the default deprecation policy
	if releaseInfo.DeprecationPolicy == nil {
		releaseInfo.DeprecationPolicy = &release.DeprecationPolicy{
			SupportedVersions: []string{version},
			ActivatesInHours:  release.DefaultDeprecationHours,
		}
	}

	b, err := json.MarshalIndent(&release.ReleaseManifest{
		Release: releaseInfo,
	}, "", "  ") // indentation: two spaces
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func parseDeprecationPolicy(releaseNotes string) (*release.DeprecationPolicy, error) {
	parts := strings.Split(releaseNotes, beginReleaseConfiguration)
	if len(parts) != 2 {
		return nil, nil
	}
	parts = strings.Split(parts[1], endReleaseConfiguration)
	if len(parts) != 2 {
		return nil, nil
	}
	policyStr := parts[0]
	if len(policyStr) == 0 {
		return nil, nil
	}
	var policy struct {
		DeprecationPolicy deprecationPolicy `yaml:"deprecationPolicy"`
	}
	if err := yaml.Unmarshal([]byte(policyStr), &policy); err != nil {
		return nil, err
	}
	return (*release.DeprecationPolicy)(&policy.DeprecationPolicy), nil
}
