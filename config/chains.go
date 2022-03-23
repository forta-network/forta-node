package config

const defaultBlockOffset = 0

var defaultRateLimiting = &RateLimitConfig{
	Rate:  0.347, // 30k/day
	Burst: 100,
}

// ChainSettings contains chain-specific settings.
type ChainSettings struct {
	Name                string
	ChainID             int
	Offset              int
	JsonRpcRateLimiting *RateLimitConfig
}

var allChainSettings = []ChainSettings{
	{
		Name:                "Ethereum Mainnet",
		ChainID:             1,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
	},
	{
		Name:                "BSC",
		ChainID:             56,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
	},
	{
		Name:                "Polygon",
		ChainID:             137,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
	},
	{
		Name:                "Avalanche",
		ChainID:             43114,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
	},
	{
		Name:                "Arbitrum",
		ChainID:             42161,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
	},
	{
		Name:                "Optimism",
		ChainID:             10,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
	},
}

// GetChainSettings returns the settings for the chain.
func GetChainSettings(chainID int) *ChainSettings {
	for _, settings := range allChainSettings {
		if settings.ChainID == chainID {
			return &settings
		}
	}
	return &ChainSettings{
		Name:                "Unknown chain",
		ChainID:             chainID,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
	}
}

// GetBlockOffset returns the block offset for a chain.
func GetBlockOffset(chainID int) int {
	return GetChainSettings(chainID).Offset
}
