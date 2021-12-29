package publisher

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	mock_clients "github.com/forta-protocol/forta-node/clients/mocks"
	"github.com/forta-protocol/forta-node/protocol"
	mock_publisher "github.com/forta-protocol/forta-node/services/publisher/mocks"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testKeyJSON       = `{"address":"ece619706cad9a43c33af9a31f1fe19f2db7cc29","crypto":{"cipher":"aes-128-ctr","ciphertext":"bef1bf4686b5576f8fd5d294a205a20c2392cd4129c5bd4ba4439d76939fea8f","cipherparams":{"iv":"c54bb09761b391efac79c942301e6a8a"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"4a72e5f0aa6e710f1e7ada8860ceb4483d29ba6e7997757c359a4f68e5d7bf3d"},"mac":"42b1c375a36ffd41e0a219edfa37b803300d1b14b951c30e1816d88c2ba31af2"},"id":"fbdc3503-010a-49e2-b950-17c326fba447","version":3}`
	testKeyPassphrase = "123"
)

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

// Suite is a test suite to test the publisher implementation.
type Suite struct {
	r *require.Assertions

	contract        *mock_publisher.MockAlertsContract
	ipfs            *mock_publisher.MockIPFS
	testAlertLogger *mock_publisher.MockTestAlertLogger
	msgClient       *mock_clients.MockMessageClient

	publisher *Publisher

	suite.Suite
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())

	s.contract = mock_publisher.NewMockAlertsContract(gomock.NewController(s.T()))
	s.ipfs = mock_publisher.NewMockIPFS(gomock.NewController(s.T()))
	s.testAlertLogger = mock_publisher.NewMockTestAlertLogger(gomock.NewController(s.T()))
	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))

	key, err := keystore.DecryptKey([]byte(testKeyJSON), testKeyPassphrase)
	s.r.NoError(err)

	cfg := PublisherConfig{
		ChainID: 1234,
		Key:     key,
	}

	s.publisher = &Publisher{
		ctx:               context.Background(),
		cfg:               cfg,
		contract:          s.contract,
		ipfs:              s.ipfs,
		testAlertLogger:   s.testAlertLogger,
		metricsAggregator: NewMetricsAggregator(),
		messageClient:     s.msgClient,

		skipEmpty:     cfg.PublisherConfig.Batch.SkipEmpty,
		skipPublish:   cfg.PublisherConfig.SkipPublish,
		batchInterval: time.Second * 1,
		batchLimit:    1,
		notifCh:       make(chan *protocol.NotifyRequest, defaultBatchLimit),
		batchCh:       make(chan *protocol.SignedAlertBatch, defaultBatchBufferSize),
	}
	// HACK: Make goroutines work by simulating the attachment message of the first agent.
	s.publisher.handleReady(nil)
}

func (s *Suite) TestPublisher() {

}
