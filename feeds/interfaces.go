package feeds

import "OpenZeppelin/fortify-node/domain"

// BlockFeed is a subscribable feed of blocks.
type BlockFeed interface {
	Start()
	Subscribe(handler func(evt *domain.BlockEvent) error) <-chan error
}

// TransactionFeed is a subscribable feed of transactions.
type TransactionFeed interface {
	ForEachTransaction(blockHandler func(evt *domain.BlockEvent) error, txHandler func(evt *domain.TransactionEvent) error) error
}
