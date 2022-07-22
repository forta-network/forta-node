package config

const defaultBlockOffset = 0

var defaultRateLimiting = &RateLimitConfig{
	Rate:  50, // 0.347, // 30k/day
	Burst: 50, // 100,
}

const defaultInspectionInterval = 100

// ChainSettings contains chain-specific settings.
type ChainSettings struct {
	Name                string
	ChainID             int
	EnableTrace         bool
	Offset              int
	JsonRpcRateLimiting *RateLimitConfig
	InspectionInterval  int // in block number
}

var allChainSettings = []ChainSettings{
	{
		Name:                "Ethereum Mainnet",
		ChainID:             1,
		EnableTrace:         true,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
		InspectionInterval:  50,
	},
	{
		Name:                "BSC",
		ChainID:             56,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
		InspectionInterval:  250,
	},
	{
		Name:                "Polygon",
		ChainID:             137,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
		InspectionInterval:  350,
	},
	{
		Name:                "Avalanche",
		ChainID:             43114,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
		InspectionInterval:  1000,
	},
	{
		Name:                "Arbitrum",
		ChainID:             42161,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
		InspectionInterval:  1000,
	},
	{
		Name:                "Optimism",
		ChainID:             10,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
		InspectionInterval:  5000,
	},
	{
		Name:                "Fantom",
		ChainID:             250,
		EnableTrace:         true,
		Offset:              defaultBlockOffset,
		JsonRpcRateLimiting: defaultRateLimiting,
		InspectionInterval:  1000,
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
		InspectionInterval:  defaultInspectionInterval,
	}
}

// GetBlockOffset returns the block offset for a chain.
func GetBlockOffset(chainID int) int {
	return GetChainSettings(chainID).Offset
}
