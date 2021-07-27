package registry

import (
	"fmt"
	"math/big"
	"testing"

	"golang.org/x/sync/semaphore"

	"OpenZeppelin/fortify-node/clients/messaging"
	mock_clients "OpenZeppelin/fortify-node/clients/mocks"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/domain"
	mock_registry "OpenZeppelin/fortify-node/services/registry/mocks"
	"OpenZeppelin/fortify-node/services/registry/regtypes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testPoolIDStr         = "0x1000000000000000000000000000000000000000000000000000000000000000"
	testAgentIDStr        = "0x2000000000000000000000000000000000000000000000000000000000000000"
	testAgentRef          = "QmWacxPov5FVCyvnpXroDJ76urakzN4ckpFhhRzpsAkRek"
	testImageRef          = "bafybeide7cspdmxqjcpa3qvrayvfpiix2it4v6mjejjc22q72zbq7rm4re@sha256:cdd4ddccf5e9c740eb4144bcc68e3ea3a056789ec7453e94a6416dcfc80937a4"
	testContainerRegistry = "some.reg.io"
	testAgentLength       = 1
)

var (
	testPoolID         = common.HexToHash(testPoolIDStr)
	testAgentID        = common.HexToHash(testAgentIDStr)
	testAgentFile      = &regtypes.AgentFile{}
	testVersion1       = big.NewInt(1)
	testVersion2       = big.NewInt(2)
	testAgentLengthBig = big.NewInt(testAgentLength)
)

// TestSuite runs the test suite.
func TestSuite(t *testing.T) {
	testAgentFile.Manifest.ImageReference = testImageRef

	suite.Run(t, &Suite{})
}

// Suite is a test suite to test the tx node runner implementation.
type Suite struct {
	r *require.Assertions

	contract   *mock_registry.MockContractRegistryCaller
	ipfsClient *mock_registry.MockIPFSClient
	ethClient  *mock_registry.MockEthClient
	msgClient  *mock_clients.MockMessageClient

	service *RegistryService

	suite.Suite
}

// SetupTest sets up the test.
func (s *Suite) SetupTest() {
	s.r = require.New(s.T())
	s.contract = mock_registry.NewMockContractRegistryCaller(gomock.NewController(s.T()))
	s.ipfsClient = mock_registry.NewMockIPFSClient(gomock.NewController(s.T()))
	s.ethClient = mock_registry.NewMockEthClient(gomock.NewController(s.T()))
	s.msgClient = mock_clients.NewMockMessageClient(gomock.NewController(s.T()))
	s.service = &RegistryService{
		poolID:     testPoolID,
		msgClient:  s.msgClient,
		contract:   s.contract,
		ipfsClient: s.ipfsClient,
		ethClient:  s.ethClient,
		done:       make(chan struct{}),
		sem:        semaphore.NewWeighted(1),
	}
	s.service.cfg.Registry.ContainerRegistry = testContainerRegistry
}

type agentConfigs []*config.AgentConfig

func (ac agentConfigs) Matches(x interface{}) bool {
	acx, ok := x.([]*config.AgentConfig)
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
	return fmt.Sprintf("%+v", ([]*config.AgentConfig)(ac))
}

func eqBytes(h common.Hash) gomock.Matcher {
	return gomock.Eq(([32]byte)(h))
}

func (s *Suite) TestDifferentVersion() {
	// Given that the last known version is 1
	s.service.version = testVersion1
	// When the last version is returned as 2 at the time of checking
	s.contract.EXPECT().PoolVersion(nil, eqBytes(s.service.poolID)).Return(testVersion2, nil)
	// Then
	s.shouldUpdateAgents()

	s.NoError(s.service.publishLatestAgents())
}

func (s *Suite) shouldUpdateAgents() {
	s.ethClient.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(&domain.Block{Number: "0x1"}, nil)
	s.contract.EXPECT().AgentLength(gomock.Any(), eqBytes(testPoolID)).Return(testAgentLengthBig, nil)
	s.contract.EXPECT().AgentAt(gomock.Any(), eqBytes(testPoolID), big.NewInt(testAgentLength-1)).
		Return(testAgentID, testAgentRef, nil)
	s.ipfsClient.EXPECT().GetAgentFile(testAgentRef).Return(testAgentFile, nil)
	s.msgClient.EXPECT().Publish(messaging.SubjectAgentsVersionsLatest, (agentConfigs)([]*config.AgentConfig{
		{
			ID:    testAgentIDStr,
			Image: fmt.Sprintf("%s/%s", testContainerRegistry, testImageRef),
		},
	}))
}

func (s *Suite) TestFirstTime() {
	// Given that there is no last known version
	s.service.version = nil
	// When the last version is returned as anything
	s.contract.EXPECT().PoolVersion(nil, eqBytes(s.service.poolID)).Return(testVersion2, nil)
	// Then
	s.shouldUpdateAgents()

	s.NoError(s.service.publishLatestAgents())
}

func (s *Suite) TestSameVersion() {
	// Given that the last known version is 1
	s.service.version = testVersion1
	// When the last version is returned as the same
	s.contract.EXPECT().PoolVersion(nil, eqBytes(s.service.poolID)).Return(testVersion1, nil)
	// Then it should silently skip

	s.NoError(s.service.publishLatestAgents())
}
