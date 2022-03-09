package config

const defaultBlockConfirmationCount = 2

type chainBlockConfirmation struct {
	Name    string
	ChainID int
	Count   int
}

var blockConfirmations = []chainBlockConfirmation{
	{
		Name:    "Ethereum Mainnet",
		ChainID: 1,
		Count:   defaultBlockConfirmationCount,
	},
}

// GetBlockConfirmationCount returns the hard coded block confirmation count for chain.
func GetBlockConfirmationCount(chainID int) int {
	for _, conf := range blockConfirmations {
		if conf.ChainID == chainID {
			return conf.Count
		}
	}
	return defaultBlockConfirmationCount
}
