package config

// Release vars - injected by the compiler
var (
	CommitHash = ""
	ReleaseCid = ""
)

// BuildReleaseInfo contains release info about current build.
type BuildReleaseInfo struct {
	Commit string `json:"commit"`
	IPFS   string `json:"ipfs"`
	// Version string `json:"version"` TODO: Use this when semver is injected
}

func GetBuildReleaseInfo() (*BuildReleaseInfo, bool) {
	if len(CommitHash) == 0 {
		return nil, false
	}

	return &BuildReleaseInfo{
		Commit: CommitHash,
		IPFS:   ReleaseCid,
	}, true
}

// ReleaseManifest contains the latest info about the latest scanner version.
type ReleaseManifest struct {
	Release Release `json:"release"`
}

// Release contains release data.
type Release struct {
	Timestamp  string          `json:"timestamp"`
	Repository string          `json:"repository"`
	Commit     string          `json:"commit"`
	Services   ReleaseServices `json:"services"`
}

// ReleaseServices are the services to run for scanner node.
type ReleaseServices struct {
	Updater    string `json:"updater"`
	Supervisor string `json:"supervisor"`
}
