package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-network/forta-core-go/ens"
	"github.com/forta-network/forta-core-go/manifest"
	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/tests/e2e"
	"github.com/forta-network/forta-node/tests/e2e/ethaccounts"
	"github.com/forta-network/forta-node/tests/e2e/misccontracts/contract_mock_registry"
	"github.com/forta-network/forta-node/tests/e2e/misccontracts/contract_multicall"
	"github.com/forta-network/forta-node/tests/e2e/misccontracts/contract_transparent_upgradeable_proxy"
	"github.com/forta-network/forta-node/testutils/alertserver"
	graphql_api "github.com/forta-network/forta-node/testutils/graphql-api"
	"github.com/golang-jwt/jwt/v4"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	smallTimeout                  = time.Minute * 3
	largeTimeout                  = time.Minute * 5
	botWithDataFeeSubscription    = "1"
	botWithoutDataFeeSubscription = "2"
)

var (
	ethereumDataDir         = ".ethereum"
	ipfsDataDir             = ".ipfs"
	genesisFile             = "genesis.json"
	passwordFile            = "ethaccounts/password"
	gethKeyFile             = "ethaccounts/gethkeyfile"
	networkID               = int64(137)
	gethNodeEndpoint        = "http://localhost:8545"
	processStartWaitSeconds = 30
	txWaitSeconds           = 5
	ipfsEndpoint            = "http://localhost:5002"
	discoConfigFile         = "disco.config.yml"
	discoPort               = "1970"

	agentID         = "0x0000000000000000000000000000000000000000000000000000000000000001"
	agentIDBigInt   = utils.AgentHexToBigInt(agentID)
	scannerIDBigInt = utils.ScannerIDHexToBigInt(ethaccounts.ScannerAddress.Hex())
	// to be set in forta-agent-0x04f4b6-02b4 format
	agentContainerID string

	runnerSupervisedContainers = []string{
		"forta-updater",
		"forta-supervisor",
	}

	allServiceContainers = []string{
		"forta-updater",
		"forta-supervisor",
		"forta-json-rpc",
		"forta-scanner",
		"forta-nats",
	}

	envEmptyRunnerTrackingID = "RUNNER_TRACKING_ID="

	testAgentLocalImageName = "forta-e2e-test-agent"

	stakeAmount, _    = big.NewInt(0).SetString("500000000000000000000", 10)
	maxStakeAmount, _ = big.NewInt(0).SetString("1000000000000000000000", 10)
)

type Suite struct {
	ctx context.Context
	r   *require.Assertions

	alertServer *alertserver.AlertServer
	graphqlAPI  *graphql_api.GraphQLAPI

	ipfsClient   *ipfsapi.Shell
	ethClient    *ethclient.Client
	dockerClient clients.DockerClient

	deployer *bind.TransactOpts
	admin    *bind.TransactOpts
	scanner  *bind.TransactOpts
	operator *bind.TransactOpts

	multicallAddr        common.Address
	mockRegistryContract *contract_mock_registry.MockRegistry

	releaseManifest    *release.ReleaseManifest
	releaseManifestCid string

	agentManifest    *manifest.SignedAgentManifest
	agentManifestCid string

	fortaProcess *Process

	suite.Suite
}

func TestE2E(t *testing.T) {
	if os.Getenv("E2E_TEST") != "1" {
		t.Log("e2e testing is not enabled (skipping) - enable with E2E_TEST=1 env var")
		return
	}

	s := &Suite{
		ctx: context.Background(),
		r:   require.New(t),
	}
	dockerClient, err := docker.NewDockerClient("")
	s.r.NoError(err)
	s.dockerClient = dockerClient

	s.ipfsClient = ipfsapi.NewShell(ipfsEndpoint)
	s.ensureAvailability("ipfs", func() error {
		_, err := s.ipfsClient.FilesLs(s.ctx, "/")
		if err != nil {
			return err
		}
		return nil
	})

	s.ensureAvailability("disco", func() error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/v2/", discoPort))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return fmt.Errorf("disco responded with status '%d'", resp.StatusCode)
	})

	ethClient, err := ethclient.Dial(gethNodeEndpoint)
	s.r.NoError(err)
	s.ethClient = ethClient
	s.ensureAvailability("geth", func() error {
		_, err := ethClient.BlockNumber(s.ctx)
		return err
	})

	suite.Run(t, s)
}

func (s *Suite) SetupTest() {
	s.ctx = context.Background()
	s.r = require.New(s.T())

	s.deployer = bind.NewKeyedTransactor(ethaccounts.DeployerKey)
	s.admin = bind.NewKeyedTransactor(ethaccounts.AccessAdminKey)
	s.scanner = bind.NewKeyedTransactor(ethaccounts.ScannerKey)
	s.operator = bind.NewKeyedTransactor(ethaccounts.ScannerOwnerKey)

	// set runtime vars and put release to ipfs
	nodeImageRef := s.readImageRef("node")
	config.DockerSupervisorImage = nodeImageRef
	config.DockerUpdaterImage = nodeImageRef
	config.UseDockerImages = "remote"
	config.Version = "0.0.1"
	s.releaseManifest = &release.ReleaseManifest{
		Release: release.Release{
			Timestamp:  time.Now().String(),
			Repository: "https://github.com/forta-network/forta-node",
			Version:    config.Version,
			Commit:     "57f35d25384ddf3f35731c636515204b1757c6ba",
			Services: release.ReleaseServices{
				Updater:    nodeImageRef,
				Supervisor: nodeImageRef,
			},
		},
	}
	s.releaseManifestCid = s.ipfsFilesAdd("/release", s.releaseManifest)
	config.ReleaseCid = s.releaseManifestCid

	// put agent manifest to ipfs
	agentImageRef := s.readImageRef("agent")
	s.agentManifest = &manifest.SignedAgentManifest{
		Manifest: &manifest.AgentManifest{
			From:           utils.StringPtr(ethaccounts.MiscAddress.Hex()),
			Name:           utils.StringPtr("Exploiter Transaction Detector"),
			AgentID:        utils.StringPtr("Exploiter Transaction Detector"),
			AgentIDHash:    utils.StringPtr(agentID),
			Version:        utils.StringPtr("0.0.1"),
			Timestamp:      utils.StringPtr(time.Now().String()),
			ImageReference: utils.StringPtr(agentImageRef),
			Repository:     utils.StringPtr("https://github.com/forta-network/forta-node/tree/master/tests/e2e/agents/txdetectoragent"),
			ChainIDs:       []int64{networkID},
		},
	}
	s.agentManifestCid = s.ipfsFilesAdd("/agent", s.agentManifest)

	agentContainerID = config.AgentConfig{
		ID:      agentID,
		Image:   agentImageRef,
		ChainID: int(networkID),
	}.ContainerName()

	// deploy multicall contract so we can load the assignment list
	multicallAddr, err := s.deployContract(
		Deploy("Multicall", s.deployer, contract_multicall.MulticallMetaData).Construct(),
	)
	s.r.NoError(err)
	s.multicallAddr = multicallAddr

	// update the public mode node config with the multicall address
	b, err := os.ReadFile(".forta/config-template.yml")
	s.r.NoError(err)
	publicModeCfg := strings.Replace(string(b), "MULTICALL_ADDRESS", s.multicallAddr.Hex(), -1)
	s.r.NoError(os.WriteFile(".forta/config.yml", []byte(publicModeCfg), 0777))

	// deploy mock contract with release and bot info
	mockRegistryAddr, err := s.deployContract(
		Deploy("MockRegistry", s.deployer, contract_mock_registry.MockRegistryMetaData).
			Construct(s.releaseManifestCid, s.agentManifestCid),
	)
	s.r.NoError(err)
	s.mockRegistryContract, _ = contract_mock_registry.NewMockRegistry(mockRegistryAddr, s.ethClient)

	// point ENS to mock
	ensOverrides := map[string]string{
		ens.DispatchContract:            mockRegistryAddr.Hex(),
		ens.AgentRegistryContract:       mockRegistryAddr.Hex(),
		ens.ScannerRegistryContract:     mockRegistryAddr.Hex(),
		ens.ScannerPoolRegistryContract: mockRegistryAddr.Hex(),
		ens.ScannerNodeVersionContract:  mockRegistryAddr.Hex(),
		ens.StakingContract:             mockRegistryAddr.Hex(),
		ens.RewardsContract:             mockRegistryAddr.Hex(),
		ens.StakeAllocatorContract:      mockRegistryAddr.Hex(),
		ens.MigrationContract:           mockRegistryAddr.Hex(),
	}
	b, _ = json.MarshalIndent(ensOverrides, "", "  ")
	s.r.NoError(ioutil.WriteFile(".forta/ens-override.json", b, 0644))
	s.r.NoError(ioutil.WriteFile(".forta-local/ens-override.json", b, 0644))

	// start the fake alert server
	s.alertServer = alertserver.New(s.ctx, 9090)
	go s.alertServer.Start()

	// start fake graphql api
	s.graphqlAPI = graphql_api.NewWithAuthMiddleware(
		s.ctx, e2e.DefaultMockGraphqlAPIPort, func(handlerFunc http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				// extract bot id information
				authHeader := r.Header.Get("Authorization")
				token := strings.Split(authHeader, "Bearer ")[1]
				t, err := security.VerifyScannerJWT(token)
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				c := t.Token.Claims.(jwt.MapClaims)
				botID := c["bot-id"]

				// only allow requests from the bot with subscription
				if botID == botWithDataFeeSubscription {
					handlerFunc.ServeHTTP(w, r)
					return
				}

				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		},
	)

	go s.graphqlAPI.Start()
}

func attachCmdOutput(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
}

func (s *Suite) runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	attachCmdOutput(cmd)
	s.r.NoError(cmd.Run())
}

func (s *Suite) runCmdSilent(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	s.r.NoError(cmd.Run())
}

func (s *Suite) ensureTx(name string, tx *types.Transaction) {
	for i := 0; i < txWaitSeconds*5; i++ {
		receipt, err := s.ethClient.TransactionReceipt(s.ctx, tx.Hash())
		if err == nil {
			s.r.Equal(tx.Hash().Hex(), receipt.TxHash.Hex())
			s.T().Logf("%s - mined: %s", name, tx.Hash())
			return
		}
		time.Sleep(time.Millisecond * 200)
	}
	time.Sleep(time.Second) // hard delay
	s.r.FailNowf("failed to mine tx", "%s: %s", name, tx.Hash())
}

type Deployment struct {
	Name            string
	Auth            *bind.TransactOpts
	ContractMeta    *bind.MetaData
	ConstructorArgs []interface{}
	InitializeArgs  []interface{}
}

func Deploy(name string, auth *bind.TransactOpts, meta *bind.MetaData) *Deployment {
	return &Deployment{
		Name:         name,
		Auth:         auth,
		ContractMeta: meta,
	}
}

func (dep *Deployment) Construct(args ...interface{}) *Deployment {
	dep.ConstructorArgs = append(dep.ConstructorArgs, args...)
	return dep
}

func (dep *Deployment) Init(args ...interface{}) *Deployment {
	dep.InitializeArgs = append(dep.InitializeArgs, args...)
	return dep
}

func (s *Suite) deployContract(dep *Deployment) (common.Address, error) {
	abi, bin := getAbiAndBin(dep.ContractMeta)
	implAddr, tx, _, err := bind.DeployContract(dep.Auth, *abi, common.FromHex(bin), s.ethClient, dep.ConstructorArgs...)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy logic contract: %v", err)
	}
	s.ensureTx(fmt.Sprintf("%s deployment", dep.Name), tx)
	return implAddr, nil
}

func (s *Suite) deployContractWithProxy(dep *Deployment) (common.Address, error) {
	implAddr, err := s.deployContract(dep)
	if err != nil {
		return common.Address{}, err
	}

	abi, _ := getAbiAndBin(dep.ContractMeta)
	initCallData, err := abi.Pack("initialize", dep.InitializeArgs...)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to construct init call data: %v", err)
	}

	proxyAddress, tx, _, err := contract_transparent_upgradeable_proxy.DeployTransparentUpgradeableProxy(
		dep.Auth, s.ethClient, implAddr, ethaccounts.ProxyAdminAddress, initCallData,
	)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy proxy: %v", err)
	}
	s.ensureTx(fmt.Sprintf("%s proxy deployment", dep.Name), tx)

	return proxyAddress, nil
}

func getAbiAndBin(metadata *bind.MetaData) (*abi.ABI, string) {
	parsed, _ := metadata.GetAbi()
	return parsed, metadata.Bin
}

func (s *Suite) readImageRef(name string) string {
	imageRefB, err := ioutil.ReadFile(fmt.Sprintf(".imagerefs/%s", name))
	s.r.NoError(err)
	imageRefB = []byte(strings.TrimSpace(string(imageRefB)))
	s.r.NotEmpty(imageRefB)
	return string(imageRefB)
}

func (s *Suite) ipfsFilesAdd(path string, data interface{}) string {
	b, err := json.Marshal(data)
	s.r.NoError(err)
	s.ipfsClient.FilesRm(s.ctx, path, true)
	err = s.ipfsClient.FilesWrite(s.ctx, path, bytes.NewBuffer(b), ipfsapi.FilesWrite.Create(true))
	s.r.NoError(err)
	stat, err := s.ipfsClient.FilesStat(s.ctx, path)
	s.r.NoError(err)
	return stat.Hash
}

func (s *Suite) ensureAvailability(name string, check func() error) {
	var err error
	for i := 0; i < processStartWaitSeconds*2; i++ {
		time.Sleep(time.Millisecond * 500)
		if err = check(); err == nil {
			return
		}
	}
	s.r.FailNowf("", "failed to ensure '%s' start: %v", name, err)
}

func (s *Suite) TearDownTest() {
	s.fortaProcess = nil
	if s.alertServer != nil {
		s.alertServer.Close()
	}
	if s.graphqlAPI != nil {
		s.graphqlAPI.Close()
	}
}

func (s *Suite) tearDownProcess(process *os.Process) {
	process.Signal(syscall.SIGINT)
	process.Wait()
}

type Process struct {
	stderr *bytes.Buffer
	stdout *bytes.Buffer
	*os.Process
}

type wrappedBuffer struct {
	w   io.Writer
	buf *bytes.Buffer
}

func (wb *wrappedBuffer) Write(b []byte) (int, error) {
	wb.buf.Write(b)
	return wb.w.Write(b)
}

func (process *Process) HasOutput(s string) bool {
	return strings.Contains(process.stdout.String(), s) || strings.Contains(process.stderr.String(), s)
}

func (process *Process) GetOutput() string {
	return process.stdout.String()
}

func (s *Suite) forta(fortaDir string, args ...string) {
	dir, err := os.Getwd()
	s.r.NoError(err)

	if fortaDir == "" {
		fortaDir = ".forta"
	}
	fullFortaDir := path.Join(dir, fortaDir)

	args = append([]string{
		"./forta-test",
	}, args...)
	cmdForta := exec.Command(args[0], args[1:]...)
	cmdForta.Env = append(cmdForta.Env,
		fmt.Sprintf("FORTA_DIR=%s", fullFortaDir),
		"FORTA_PASSPHRASE=0",
	)
	var (
		stderrBuf bytes.Buffer
		stdoutBuf bytes.Buffer
	)
	cmdForta.Stderr = &wrappedBuffer{w: os.Stderr, buf: &stderrBuf}
	cmdForta.Stdout = &wrappedBuffer{w: os.Stdout, buf: &stdoutBuf}

	s.r.NoError(cmdForta.Start())
	s.T().Log("forta cmd started")
	s.fortaProcess = &Process{
		stderr:  &stderrBuf,
		stdout:  &stdoutBuf,
		Process: cmdForta.Process,
	}
}

func (s *Suite) startForta(register ...bool) {
	s.forta("", "run")
	s.expectUpIn(largeTimeout, allServiceContainers...)
}

func (s *Suite) stopForta() {
	s.r.NoError(s.fortaProcess.Signal(syscall.SIGINT))
	// s.expectDownIn(largeTimeout, allServiceContainers...)
	_, err := s.fortaProcess.Wait()
	s.r.NoError(err)
}

func (s *Suite) expectIn(timeout time.Duration, conditionFunc func() bool) {
	start := time.Now()
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		if time.Since(start) > timeout {
			s.r.FailNow("expectIn() timed out")
			return
		}
		if ok := conditionFunc(); ok {
			return
		}
	}
}

func (s *Suite) expectUpIn(timeout time.Duration, containerNames ...string) {
	s.expectIn(timeout, func() bool {
		containers, err := s.dockerClient.GetContainers(s.ctx)
		s.r.NoError(err)
		for _, containerName := range containerNames {
			container, ok := containers.ContainsAny(containerName)
			if !ok {
				return false
			}
			if container.State != "running" {
				return false
			}
		}
		return true
	})
}

func (s *Suite) expectDownIn(timeout time.Duration, containerNames ...string) {
	s.expectIn(timeout, func() bool {
		containers, err := s.dockerClient.GetContainers(s.ctx)
		s.r.NoError(err)
		for _, containerName := range containerNames {
			container, ok := containers.FindByName(containerName)
			if !ok {
				continue
			}
			if ok && container.State != "exited" {
				return false
			}
		}
		return true
	})
}

func (s *Suite) sendExploiterTx() {
	gasPrice, err := s.ethClient.SuggestGasPrice(s.ctx)
	s.r.NoError(err)
	nonce, err := s.ethClient.PendingNonceAt(s.ctx, ethaccounts.ExploiterAddress)
	s.r.NoError(err)
	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &ethaccounts.ExploiterAddress,
		Value:    big.NewInt(1),
		GasPrice: gasPrice,
		Gas:      100000, // 100k
	}
	tx, err := types.SignNewTx(ethaccounts.ExploiterKey, types.HomesteadSigner{}, txData)
	s.r.NoError(err)

	s.r.NoError(s.ethClient.SendTransaction(s.ctx, tx))
	s.ensureTx("Exploiter account sending 1 Wei to itself", tx)
}
