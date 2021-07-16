package registry

import (
	"OpenZeppelin/fortify-node/clients/messaging"
	mock_clients "OpenZeppelin/fortify-node/clients/mocks"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/domain"
	mock_feeds "OpenZeppelin/fortify-node/feeds/mocks"
	mock_registry "OpenZeppelin/fortify-node/services/registry/mocks"
	"OpenZeppelin/fortify-node/services/registry/regtypes"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testPoolIDStr  = "0x1000000000000000000000000000000000000000000000000000000000000000"
	testAgentIDStr = "0x2000000000000000000000000000000000000000000000000000000000000000"
	testAgentRef   = "QmWacxPov5FVCyvnpXroDJ76urakzN4ckpFhhRzpsAkRek"
	testImageRef   = "bafybeide7cspdmxqjcpa3qvrayvfpiix2it4v6mjejjc22q72zbq7rm4re@sha256:cdd4ddccf5e9c740eb4144bcc68e3ea3a056789ec7453e94a6416dcfc80937a4"
)

var (
	testPoolID  = common.HexToHash(testPoolIDStr)
	testAgentID = common.HexToHash(testAgentIDStr)
	testLog     = &types.Log{}
	testTx      = &domain.TransactionEvent{
		Receipt: &domain.TransactionReceipt{
			Logs: []domain.LogEntry{
				{},
			},
		},
	}
	testAgentFile = &regtypes.AgentFile{}
)

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	testAgentFile.Manifest.ImageReference = testImageRef

	suite.Run(t, &Suite{})
}

// Suite is a test suite to test the tx node runner implementation.
type Suite struct {
	r *require.Assertions

	txFeed      *mock_feeds.MockTransactionFeed
	contract    *mock_registry.MockContractRegistryCaller
	logUnpacker *mock_registry.MockLogUnpacker
	ipfsClient  *mock_registry.MockIPFSClient
	msgClient   *mock_clients.MockMessageClient

	service *RegistryService

	suite.Suite
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())
	s.txFeed = mock_feeds.NewMockTransactionFeed(gomock.NewController(s.T()))
	s.contract = mock_registry.NewMockContractRegistryCaller(gomock.NewController(s.T()))
	s.logUnpacker = mock_registry.NewMockLogUnpacker(gomock.NewController(s.T()))
	s.ipfsClient = mock_registry.NewMockIPFSClient(gomock.NewController(s.T()))
	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))
	s.service = &RegistryService{
		poolID:       common.HexToHash(testPoolIDStr),
		msgClient:    s.msgClient,
		txFeed:       s.txFeed,
		contract:     s.contract,
		logUnpacker:  s.logUnpacker,
		ipfsClient:   s.ipfsClient,
		agentUpdates: make(chan *agentUpdate, 100),
		done:         make(chan struct{}),
	}
	s.txFeed.EXPECT().ForEachTransaction(nil, gomock.Any()).AnyTimes()
	s.contract.EXPECT().AgentLength(nil, gomock.Any()).Return(big.NewInt(0), nil)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsVersionsLatest, (agentConfigs)([]config.AgentConfig{}))
	s.r.NoError(s.service.start())
}

type agentConfigs []config.AgentConfig

func (ac agentConfigs) Matches(x interface{}) bool {
	acx, ok := x.([]config.AgentConfig)
	if !ok {
		return false
	}

	if len(ac) != len(acx) {
		return false
	}

	for i, agent1 := range ac {
		agent2 := acx[i]
		if !(agent1.ID == agent2.ID && agent1.Image == agent2.Image) {
			return false
		}
	}
	return true
}

func (ac agentConfigs) String() string {
	return fmt.Sprintf("%+v", ([]config.AgentConfig)(ac))
}

func (s *Suite) TestAgentAddRemove() {
	s.service.agentUpdatesWg.Add(2)

	// Add agent

	s.logUnpacker.EXPECT().UnpackAgentRegistryAgentAdded(testLog).Return(&contracts.AgentRegistryAgentAdded{
		PoolId:  common.HexToHash(testPoolIDStr),
		AgentId: common.HexToHash(testAgentIDStr),
		Ref:     testAgentRef,
	}, nil)
	s.ipfsClient.EXPECT().GetAgentFile(testAgentRef).Return(testAgentFile, nil)
	// Final state: One agent
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsVersionsLatest, (agentConfigs)([]config.AgentConfig{
		{
			ID:    testAgentIDStr,
			Image: testImageRef,
		},
	}))

	update, agentID, ref, err := s.service.detectAgentEvent(testTx)
	s.r.NoError(err)
	s.r.NoError(s.service.sendAgentUpdate(update, agentID, ref))
	s.service.agentUpdatesWg.Done()

	// Remove agent

	s.logUnpacker.EXPECT().UnpackAgentRegistryAgentAdded(testLog).Return(nil, errors.New("some error"))
	s.logUnpacker.EXPECT().UnpackAgentRegistryAgentUpdated(testLog).Return(nil, errors.New("some error"))
	s.logUnpacker.EXPECT().UnpackAgentRegistryAgentRemoved(testLog).Return(&contracts.AgentRegistryAgentRemoved{
		PoolId:  common.HexToHash(testPoolIDStr),
		AgentId: common.HexToHash(testAgentIDStr),
	}, nil)
	s.ipfsClient.EXPECT().GetAgentFile(testAgentRef).Return(testAgentFile, nil)
	// Final state: No agents
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsVersionsLatest, (agentConfigs)([]config.AgentConfig{}))

	update, agentID, ref, err = s.service.detectAgentEvent(testTx)
	s.r.NoError(err)
	s.r.NoError(s.service.sendAgentUpdate(update, agentID, ref))
	close(s.service.agentUpdates)
	s.service.agentUpdatesWg.Done()
	<-s.service.done
}
