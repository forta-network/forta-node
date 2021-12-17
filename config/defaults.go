package config

const (
	DefaultLocalAgentsFileName = "local-agents.json"
	DefaultKeysDirName         = ".keys"
	DefaultNatsPort            = "4222"
	DefaultFortaNodeBinaryPath = "/forta-node" // the path for the common binary in the container image
)

// ENS contains the default names.
type ENS struct {
	Dispatch string
	Alerts   string
	Agents   string
}

// GetENSNames returns the default ENS names.
func GetENSNames() *ENS {
	return &ENS{
		Dispatch: "dispatch.forta.eth",
		Alerts:   "alerts.forta.eth",
		Agents:   "agents.registries.forta.eth",
	}
}
