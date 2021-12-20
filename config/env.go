package config

const (
	EnvFortaDir      = "FORTA_DIR"
	EnvJsonRpcHost   = "JSON_RPC_HOST"
	EnvJsonRpcPort   = "JSON_RPC_PORT"
	EnvAgentGrpcPort = "AGENT_GRPC_PORT"
)

// EnvDefaults contain default values for one env.
type EnvDefaults struct {
	DiscoSubdomain string
}

// GetEnvDefaults returns the default values for an env.
func GetEnvDefaults(development bool) EnvDefaults {
	if development {
		return EnvDefaults{
			DiscoSubdomain: "disco-dev",
		}
	}
	return EnvDefaults{
		DiscoSubdomain: "disco",
	}
}
