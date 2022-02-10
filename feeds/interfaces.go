package feeds

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/domain"
)

// BlockFeed is a subscribable feed of blocks.
type BlockFeed interface {
	Start()
	StartRange(start int64, end int64, rate int64)
	IsStarted() bool
	Subscribe(handler func(evt *domain.BlockEvent) error) <-chan error
	health.Reporter
}

// TransactionFeed is a subscribable feed of transactions.
type TransactionFeed interface {
	ForEachTransaction(blockHandler func(evt *domain.BlockEvent) error, txHandler func(evt *domain.TransactionEvent) error) error
}

// LogFeed is a feed of logs
type LogFeed interface {
	ForEachLog(handler func(blk *domain.Block, logEntry types.Log) error, finishBlockHandler func(blk *domain.Block) error) error
}
