package publisher

import (
	"context"
	"io"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mock_clients "github.com/forta-protocol/forta-node/clients/mocks"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/protocol"
	mock_publisher "github.com/forta-protocol/forta-node/services/publisher/mocks"
	"github.com/golang/mock/gomock"
	ipfsapi "github.com/ipfs/go-ipfs-api"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testKeyJSON       = `{"address":"ece619706cad9a43c33af9a31f1fe19f2db7cc29","crypto":{"cipher":"aes-128-ctr","ciphertext":"bef1bf4686b5576f8fd5d294a205a20c2392cd4129c5bd4ba4439d76939fea8f","cipherparams":{"iv":"c54bb09761b391efac79c942301e6a8a"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"4a72e5f0aa6e710f1e7ada8860ceb4483d29ba6e7997757c359a4f68e5d7bf3d"},"mac":"42b1c375a36ffd41e0a219edfa37b803300d1b14b951c30e1816d88c2ba31af2"},"id":"fbdc3503-010a-49e2-b950-17c326fba447","version":3}`
	testKeyPassphrase = "123"
	testIPFSRef       = "QmQzYmb815G6B7MFcwPGzS71atB1toN83HgwF5xxaTZ5n8"
)

// TestPublisherSuite runs the test suite.
func TestPublisherSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

// Suite is a test suite to test the publisher implementation.
type Suite struct {
	r *require.Assertions

	contract        *mock_publisher.MockAlertsContract
	ipfs            *ipfsClientStub
	testAlertLogger *mock_publisher.MockTestAlertLogger
	msgClient       *mock_clients.MockMessageClient

	publisher *Publisher

	suite.Suite
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())

	s.contract = mock_publisher.NewMockAlertsContract(gomock.NewController(s.T()))
	s.ipfs = &ipfsClientStub{payloadCh: make(chan []byte)}
	s.testAlertLogger = mock_publisher.NewMockTestAlertLogger(gomock.NewController(s.T()))
	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))

	key, err := keystore.DecryptKey([]byte(testKeyJSON), testKeyPassphrase)
	s.r.NoError(err)

	cfg := PublisherConfig{
		ChainID: 1234,
		Key:     key,
		PublisherConfig: config.PublisherConfig{
			Batch: config.BatchConfig{
				SkipEmpty: true,
			},
		},
	}

	s.publisher = &Publisher{
		ctx:               context.Background(),
		cfg:               cfg,
		contract:          s.contract,
		ipfs:              s.ipfs,
		testAlertLogger:   s.testAlertLogger,
		metricsAggregator: NewMetricsAggregator(),
		messageClient:     s.msgClient,

		skipEmpty:     true,
		skipPublish:   false,
		batchInterval: time.Second * 1,
		batchLimit:    1,
		notifCh:       make(chan *protocol.NotifyRequest, defaultBatchLimit),
		batchCh:       make(chan *protocol.SignedAlertBatch, defaultBatchBufferSize),
	}

	s.msgClient.EXPECT().Subscribe(gomock.Any(), gomock.Any()).AnyTimes()

	// HACK: Make goroutines work by simulating the attachment message of the first agent.
	s.publisher.handleReady(nil)
}

type ipfsClientStub struct {
	payloadCh chan []byte
}

func (client *ipfsClientStub) Add(r io.Reader, options ...ipfsapi.AddOpts) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}
	client.payloadCh <- b
	return testIPFSRef, nil
}

func (s *Suite) TestPublisher() {
	req := &protocol.NotifyRequest{
		AgentInfo: &protocol.AgentInfo{},
		SignedAlert: &protocol.SignedAlert{
			Alert: &protocol.Alert{
				Agent: &protocol.AgentInfo{},
				Finding: &protocol.Finding{
					Severity: protocol.Finding_INFO,
				},
			},
		},
		EvalTxRequest: &protocol.EvaluateTxRequest{
			RequestId: "1234",
			Event: &protocol.TransactionEvent{
				Block: &protocol.TransactionEvent_EthBlock{
					BlockNumber: "0x1",
				},
				Receipt: &protocol.TransactionEvent_EthReceipt{
					TransactionHash: "0xabc",
				},
			},
		},
	}

	s.contract.EXPECT().AddAlertBatch(
		gomock.Any(), gomock.Any(), gomock.Any(), big.NewInt(1),
		big.NewInt(int64(protocol.Finding_INFO)), testIPFSRef,
	).Return(makeTestTx(), nil)

	s.publisher.Notify(s.publisher.ctx, req)

	batchPayload := <-s.ipfs.payloadCh
	s.T().Log(string(batchPayload))

	s.r.Contains(string(batchPayload), req.EvalTxRequest.RequestId)
}

func makeTestTx() *types.Transaction {
	return types.NewTransaction(0, common.HexToAddress("0xcB548DD68835F12244F702a316c5eDA106f7708C"), big.NewInt(0), 21000, big.NewInt(30000), nil)
}
