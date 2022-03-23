package config

const defaultBlockOffset = 0

var defaultRateLimiting = &RateLimitConfig{
	Rate:  0.347, // 30k/day
	Burst: 100,
}

// ChainSettings contains chain-specific settings.
type ChainSettings struct {
	Name         string
	ChainID      int
	Offset       int
	RateLimiting *RateLimitConfig
}

var allChainSettings = []ChainSettings{
	{
		Name:         "Ethereum Mainnet",
		ChainID:      1,
		Offset:       defaultBlockOffset,
		RateLimiting: defaultRateLimiting,
	},
	{
		Name:         "BSC",
		ChainID:      56,
		Offset:       defaultBlockOffset,
		RateLimiting: defaultRateLimiting,
	},
	{
		Name:         "Polygon",
		ChainID:      137,
		Offset:       defaultBlockOffset,
		RateLimiting: defaultRateLimiting,
	},
	{
		Name:         "Avalanche",
		ChainID:      43114,
		Offset:       defaultBlockOffset,
		RateLimiting: defaultRateLimiting,
	},
	{
		Name:         "Arbitrum",
		ChainID:      42161,
		Offset:       defaultBlockOffset,
		RateLimiting: defaultRateLimiting,
	},
	{
		Name:         "Optimism",
		ChainID:      10,
		Offset:       defaultBlockOffset,
		RateLimiting: defaultRateLimiting,
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
		Name:         "Unknown chain",
		ChainID:      chainID,
		Offset:       defaultBlockOffset,
		RateLimiting: defaultRateLimiting,
	}
}

// GetBlockOffset returns the block offset for a chain.
func GetBlockOffset(chainID int) int {
	return GetChainSettings(chainID).Offset
}
