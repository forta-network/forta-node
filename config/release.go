package config

// ReleaseManifest contains the latest info about the latest scanner version.
type ReleaseManifest struct {
	Metadata ReleaseMetadata `json:"metadata"`
}

// ReleaseMetadata contains release data.
type ReleaseMetadata struct {
	Version   string          `json:"version"`
	Timestamp string          `json:"timestamp"`
	Notes     string          `json:"notes"`
	Services  ReleaseServices `json:"services"`
}

// ReleaseServices are the services to run for scanner node.
type ReleaseServices struct {
	Updater    string `json:"updater"`
	Supervisor string `json:"supervisor"`
}
