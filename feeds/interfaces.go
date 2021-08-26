package feeds

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/forta-network/forta-node/domain"
)

// BlockFeed is a subscribable feed of blocks.
type BlockFeed interface {
	Start()
	Subscribe(handler func(evt *domain.BlockEvent) error) <-chan error
}

// TransactionFeed is a subscribable feed of transactions.
type TransactionFeed interface {
	ForEachTransaction(blockHandler func(evt *domain.BlockEvent) error, txHandler func(evt *domain.TransactionEvent) error) error
}

// LogFeed is a feed of logs
type LogFeed interface {
	ForEachLog(handler func(log types.Log) error) error
}
