package config

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
