package config

const (
	EnvHostFortaDir = "HOST_FORTA_DIR" // for retrieving forta dir path on the host os
	EnvDevelopment  = "FORTA_DEVELOPMENT"
	EnvReleaseInfo  = "FORTA_RELEASE_INFO"

	// Agent env vars
	EnvJWTProviderHost       = "FORTA_JWT_PROVIDER_HOST"
	EnvJWTProviderPort       = "FORTA_JWT_PROVIDER_PORT"
	EnvPublicAPIProxyHost    = "FORTA_PUBLIC_API_PROXY_HOST"
	EnvPublicAPIProxyPort    = "FORTA_PUBLIC_API_PROXY_PORT"
	EnvAgentGrpcPort         = "AGENT_GRPC_PORT"
	EnvFortaHealthCheckPort  = "FORTA_HEALTH_CHECK_PORT"
	EnvFortaBotID            = "FORTA_BOT_ID"
	EnvFortaBotOwner         = "FORTA_BOT_OWNER"
	EnvFortaChainID          = "FORTA_CHAIN_ID"
	EnvFortaShardID          = "FORTA_SHARD_ID"
	EnvFortaShardCount       = "FORTA_SHARD_COUNT"
	EnvFortaTokenExchangeURL = "FORTA_TOKEN_EXCHANGE_URL"

	EnvJsonRpcHost           = "JSON_RPC_HOST"
	EnvJsonRpcPort           = "JSON_RPC_PORT"
	EnvCacheJsonRpcCachePort = "CACHE_JSON_RPC_PORT"
	EnvCacheRequestTimeout   = "CACHE_REQUEST_TIMEOUT"
	EnvCacheRequestInterval  = "CACHE_REQUEST_INTERVAL"
	EnvCacheSupportedChains  = "CACHE_SUPPORTED_CHAINS"
)
