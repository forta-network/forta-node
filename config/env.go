package config

const (
	EnvHostFortaDir = "HOST_FORTA_DIR" // for retrieving forta dir path on the host os
	EnvDevelopment  = "FORTA_DEVELOPMENT"
	EnvReleaseInfo  = "FORTA_RELEASE_INFO"

	// Agent env vars
	EnvJsonRpcHost        = "JSON_RPC_HOST"
	EnvJsonRpcPort        = "JSON_RPC_PORT"
	EnvJWTProviderHost    = "FORTA_JWT_PROVIDER_HOST"
	EnvJWTProviderPort    = "FORTA_JWT_PROVIDER_PORT"
	EnvPublicAPIProxyHost = "FORTA_PUBLIC_API_PROXY_HOST"
	EnvPublicAPIProxyPort = "FORTA_PUBLIC_API_PROXY_PORT"
	EnvAgentGrpcPort      = "AGENT_GRPC_PORT"
	EnvFortaBotID         = "FORTA_BOT_ID"
	EnvFortaBotOwner      = "FORTA_BOT_OWNER"
	EnvFortaChainID       = "FORTA_CHAIN_ID"
	EnvFortaShardID       = "FORTA_SHARD_ID"
	EnvFortaShardCount    = "FORTA_SHARD_COUNT"
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
