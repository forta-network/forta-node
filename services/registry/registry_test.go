package registry

import (
	"fmt"
	"math/big"
	"testing"

	"golang.org/x/sync/semaphore"

	"github.com/forta-network/forta-node/clients/messaging"
	mock_clients "github.com/forta-network/forta-node/clients/mocks"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/domain"
	mock_registry "github.com/forta-network/forta-node/services/registry/mocks"
	"github.com/forta-network/forta-node/services/registry/regtypes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testScannerAddressStr = "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"
	testAgentIDStr        = "0x2000000000000000000000000000000000000000000000000000000000000000"
	testAgentRef          = "QmWacxPov5FVCyvnpXroDJ76urakzN4ckpFhhRzpsAkRek"
	testImageRef          = "bafybeide7cspdmxqjcpa3qvrayvfpiix2it4v6mjejjc22q72zbq7rm4re@sha256:cdd4ddccf5e9c740eb4144bcc68e3ea3a056789ec7453e94a6416dcfc80937a4"
	testContainerRegistry = "some.reg.io"
	testAgentLength       = 1
)

var (
	testScannerAddress = common.HexToAddress(testScannerAddressStr)
	testAgentID        = common.HexToHash(testAgentIDStr)
	testAgentFile      = &regtypes.AgentFile{}
	testVersion1       = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	testVersion2       = [32]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
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
		scannerAddress: testScannerAddress,
		msgClient:      s.msgClient,
		contract:       s.contract,
		ipfsClient:     s.ipfsClient,
		ethClient:      s.ethClient,
		done:           make(chan struct{}),
		sem:            semaphore.NewWeighted(1),
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

func (s *Suite) TestDifferentVersion() {
	// Given that the last known version is 1
	s.service.version = string(testVersion1[:])
	// When the last version is returned as 2 at the time of checking
	s.contract.EXPECT().GetAgentListHash(nil, s.service.scannerAddress).Return(testVersion2, nil)
	// Then
	s.shouldUpdateAgents()

	s.NoError(s.service.publishLatestAgents())
}

func (s *Suite) shouldUpdateAgents() {
	s.ethClient.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(&domain.Block{Number: "0x1"}, nil)
	s.contract.EXPECT().AgentLength(gomock.Any(), testScannerAddress).Return(testAgentLengthBig, nil)
	s.contract.EXPECT().AgentAt(gomock.Any(), testScannerAddress, big.NewInt(testAgentLength-1)).
		Return(testAgentID, big.NewInt(0), false, testAgentRef, false, nil)
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
	s.service.version = ""
	// When the last version is returned as anything
	s.contract.EXPECT().GetAgentListHash(nil, s.service.scannerAddress).Return(testVersion2, nil)
	// Then
	s.shouldUpdateAgents()

	s.NoError(s.service.publishLatestAgents())
}

func (s *Suite) TestSameVersion() {
	// Given that the last known version is 1
	s.service.version = string(testVersion1[:])
	// When the last version is returned as the same
	s.contract.EXPECT().GetAgentListHash(nil, s.service.scannerAddress).Return(testVersion1, nil)
	// Then it should silently skip

	s.NoError(s.service.publishLatestAgents())
}
