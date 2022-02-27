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

// ReleaseSummary contains concise release info.
type ReleaseSummary struct {
	Commit  string `json:"commit,omitempty"`
	IPFS    string `json:"ipfs,omitempty"`
	Version string `json:"version,omitempty"`
}

// GetBuildReleaseSummary returns the build summary from build vars.
func GetBuildReleaseSummary() (*ReleaseSummary, bool) {
	if len(CommitHash) == 0 {
		return nil, false
	}

	return &ReleaseSummary{
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
