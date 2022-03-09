package config

const defaultBlockOffset = 0

// ChainSettings contains chain-specific settings.
type ChainSettings struct {
	Name    string
	ChainID int
	Offset  int
}

var allChainSettings = []ChainSettings{
	{
		Name:    "Ethereum Mainnet",
		ChainID: 1,
		Offset:  defaultBlockOffset,
	},
}

// GetChainSettings returns the settings for the chain.
func GetChainSettings(chainID int) (*ChainSettings, bool) {
	for _, settings := range allChainSettings {
		if settings.ChainID == chainID {
			return &settings, true
		}
	}
	return nil, false
}

// GetBlockOffset returns the block offset for a chain.
func GetBlockOffset(chainID int) int {
	settings, ok := GetChainSettings(chainID)
	if ok {
		return settings.Offset
	}
	return defaultBlockOffset
}
