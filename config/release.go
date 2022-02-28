package config

import (
	"github.com/forta-protocol/forta-core-go/release"
)

// Release vars - injected by the compiler
var (
	CommitHash = ""
	ReleaseCid = ""
	Version    = ""
)

// GetBuildReleaseSummary returns the build summary from build vars.
func GetBuildReleaseSummary() (*release.ReleaseSummary, bool) {
	if len(CommitHash) == 0 {
		return nil, false
	}

	return &release.ReleaseSummary{
		Commit:  CommitHash,
		IPFS:    ReleaseCid,
		Version: Version,
	}, true
}

// GetBuildReleaseInfo collects and returns the release info from build vars.
func GetBuildReleaseInfo() *release.ReleaseInfo {
	return &release.ReleaseInfo{
		FromBuild: true,
		IPFS:      ReleaseCid,
		Manifest: release.ReleaseManifest{
			Release: release.Release{
				Version: Version,
				Commit:  CommitHash,
			},
		},
	}
}
