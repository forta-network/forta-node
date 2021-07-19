package regtypes

// AgentFile contains agent data.
type AgentFile struct {
	Manifest struct {
		ImageReference string `json:"imageReference"`
	} `json:"manifest"`
}
