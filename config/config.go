package config

import (
	"fmt"
	"math/big"
	"os"
	"path"

	"github.com/goccy/go-json"

	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/utils"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

const EnvConfig = "FORTA_CONFIG"
const EnvFortaDir = "FORTA_DIR"
const EnvJsonRpcHost = "JSON_RPC_HOST"
const EnvJsonRpcPort = "JSON_RPC_PORT"
const EnvAgentGrpcPort = "AGENT_GRPC_PORT"
const ContainerNamePrefix = "forta"

const (
	//TODO: figure out protocol for this other than direct communication
	EnvPublisherHost = "PUBLISHER_HOST"
	// NATS node name env var
	EnvNatsHost = "NATS_HOST"
)

const (
	DefaultLocalAgentsFileName = "local-agents.json"
	DefaultKeysDirName         = ".keys"
	DefaultFortaNodeBinaryPath = "/forta-node" // the path for the common binary in the container image
)

// Docker container names
var (
	DockerScannerNodeImage = "forta-protocol/forta-node:latest"
	UseDockerImages        = "local"

	DockerSupervisorContainerName   = fmt.Sprintf("%s-supervisor", ContainerNamePrefix)
	DockerNatsContainerName         = fmt.Sprintf("%s-nats", ContainerNamePrefix)
	DockerScannerContainerName      = fmt.Sprintf("%s-scanner", ContainerNamePrefix)
	DockerJSONRPCProxyContainerName = fmt.Sprintf("%s-json-rpc", ContainerNamePrefix)
	DockerPublisherContainerName    = fmt.Sprintf("%s-publisher", ContainerNamePrefix)

	DockerNetworkName = DockerScannerContainerName

	DefaultContainerFortaDirPath        = "/.forta"
	DefaultContainerKeyDirPath          = path.Join(DefaultContainerFortaDirPath, DefaultKeysDirName)
	DefaultContainerLocalAgentsFilePath = path.Join(DefaultContainerFortaDirPath, DefaultLocalAgentsFileName)
)

// Global constant values
var (
	DefaultNatsPort    = "4222"
	DefaultIPFSGateway = "https://cloudflare-ipfs.com"
)

type AgentConfig struct {
	ID         string  `yaml:"id" json:"id"`
	Image      string  `yaml:"image" json:"image"`
	Manifest   string  `yaml:"manifest" json:"manifest"`
	IsLocal    bool    `yaml:"isLocal" json:"isLocal"`
	StartBlock *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock  *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
}

// ToAgentInfo transforms the agent config to the agent info.
func (ac AgentConfig) ToAgentInfo() *protocol.AgentInfo {
	return &protocol.AgentInfo{
		Id:        ac.ID,
		Image:     ac.Image,
		ImageHash: ac.ImageHash(),
		IsTest:    ac.IsLocal,
		Manifest:  ac.Manifest,
	}
}

func (ac AgentConfig) ImageHash() string {
	_, digest := utils.SplitImageRef(ac.Image)
	return digest
}

func (ac AgentConfig) ContainerName() string {
	_, digest := utils.SplitImageRef(ac.Image)
	return fmt.Sprintf("%s-agent-%s-%s", ContainerNamePrefix, utils.ShortenString(ac.ID, 8), utils.ShortenString(digest, 4))
}

func (ac AgentConfig) GrpcPort() string {
	return "50051"
}

type EthereumConfig struct {
	JsonRpcUrl   string            `yaml:"jsonRpcUrl" json:"jsonRpcUrl" validate:"omitempty,url"`
	WebsocketUrl string            `yaml:"websocketUrl" json:"websocketUrl" validate:"required_without=JsonRpcUrl"`
	Headers      map[string]string `yaml:"headers" json:"headers"`
}

type QueryConfig struct {
	Port      int             `yaml:"port" json:"port" validate:"max=65535"`
	PublishTo PublisherConfig `yaml:"publishTo" json:"publishTo"`
}
type ScannerConfig struct {
	ChainID          int            `yaml:"chainId" json:"chainId"`
	StartBlock       int            `yaml:"startBlock" json:"startBlock"`
	EndBlock         int            `yaml:"endBlock" json:"endBlock"`
	Ethereum         EthereumConfig `yaml:"ethereum" json:"ethereum"`
	DisableAutostart bool           `yaml:"disableAutostart" json:"disableAutostart"`
	BlockRateLimit   int            `yaml:"blockRateLimit" json:"blockRateLimit"`
}

type TraceConfig struct {
	Ethereum EthereumConfig `yaml:"ethereum" json:"ethereum"`
	Enabled  bool           `yaml:"enabled" json:"enabled"`
}

type JsonRpcProxyConfig struct {
	Ethereum EthereumConfig `yaml:"ethereum" json:"ethereum"`
}

type LogConfig struct {
	Level       string `yaml:"level" json:"level"`
	MaxLogSize  string `yaml:"maxLogSize" json:"maxLogSize"`
	MaxLogFiles int    `yaml:"maxLogFiles" json:"maxLogFiles"`
}

type RegistryConfig struct {
	Ethereum          EthereumConfig `yaml:"ethereum" json:"ethereum"`
	IPFSGateway       *string        `yaml:"ipfsGateway" json:"ipfsGateway,omitempty" validate:"omitempty,url"`
	ContractAddress   string         `yaml:"contractAddress" json:"contractAddress" validate:"eth_addr"`
	ContainerRegistry string         `yaml:"containerRegistry" json:"containerRegistry" validate:"hostname"`
	Username          string         `yaml:"username" json:"username"`
	Password          string         `yaml:"password" json:"password"`
	Disabled          bool           `yaml:"disabled" json:"disabled"` // for testing situations
}

type IPFSConfig struct {
	GatewayURL string `yaml:"gatewayUrl" json:"gatewayUrl" validate:"url"`
	Username   string `yaml:"username" json:"username"`
	Password   string `yaml:"password" json:"password"`
}

type BatchConfig struct {
	SkipEmpty       bool `yaml:"skipEmpty" json:"skipEmpty"`
	IntervalSeconds *int `yaml:"intervalSeconds" json:"intervalSeconds"`
	MaxAlerts       *int `yaml:"maxAlerts" json:"maxAlerts"`
}

type TestAlertsConfig struct {
	Disable    bool   `yaml:"disable" json:"disable"`
	WebhookURL string `yaml:"webhookUrl" json:"webhookUrl" validate:"omitempty,url"`
}

type PublisherConfig struct {
	SkipPublish     bool             `yaml:"skipPublish" json:"skipPublish"`
	Ethereum        EthereumConfig   `yaml:"ethereum" json:"ethereum"  validate:"required_unless=SkipPublish true"`
	ContractAddress string           `yaml:"contractAddress" json:"contractAddress" validate:"required_unless=SkipPublish true,omitempty,eth_addr"`
	IPFS            *IPFSConfig      `yaml:"ipfs" json:"ipfs" validate:"required_unless=SkipPublish true"`
	Batch           BatchConfig      `yaml:"batch" json:"batch"`
	TestAlerts      TestAlertsConfig `yaml:"testAlerts" json:"testAlerts"`
	GasPriceGwei    *int64           `yaml:"gasPriceGwei" json:"gasPriceGwei"`
	GasLimit        *uint64          `yaml:"gasLimit" json:"gasLimit"`
}

type ResourcesConfig struct {
	DisableAgentLimits bool    `yaml:"disableAgentLimits" json:"disableAgentLimits"`
	AgentMaxMemoryMiB  int     `yaml:"agentMaxMemoryMib" json:"agentMaxMemoryMib" validate:"omitempty,min=100"`
	AgentMaxCPUs       float64 `yaml:"agentMaxCpus" json:"agentMaxCpus" validate:"omitempty,gt=0"`
}

type ENSConfig struct {
	ContractAddress string          `yaml:"contractAddress" json:"contractAddress" validate:"omitempty,eth_addr"`
	Ethereum        *EthereumConfig `yaml:"ethereum" json:"ethereum"`
}

type Config struct {
	Development                  bool           `yaml:"-" json:"_development"`
	FortaDir                     string         `yaml:"-" json:"_fortaDir"`
	ConfigPath                   string         `yaml:"-" json:"_configPath"`
	KeyDirPath                   string         `yaml:"-" json:"_keyDirPath"`
	Passphrase                   string         `yaml:"-" json:"_passphrase"`
	LocalAgentsPath              string         `yaml:"-" json:"_localAgentsPath"`
	LocalAgents                  []*AgentConfig `yaml:"-" json:"_localAgents"`
	AgentRegistryContractAddress string         `yaml:"-" json:"_agentRegistry"`

	Registry        RegistryConfig     `yaml:"registry" json:"registry"`
	Scanner         ScannerConfig      `yaml:"scanner" json:"scanner"`
	Query           QueryConfig        `yaml:"query" json:"query"`
	Trace           TraceConfig        `yaml:"trace" json:"trace"`
	JsonRpcProxy    JsonRpcProxyConfig `yaml:"jsonRpcProxy" json:"jsonRpcProxy"`
	Log             LogConfig          `yaml:"log" json:"log"`
	ResourcesConfig ResourcesConfig    `yaml:"resources" json:"resources"`
	ENSConfig       ENSConfig          `yaml:"ens" json:"ens"`
}

func ParseBigInt(num int) *big.Int {
	var val *big.Int
	if num != 0 {
		val = big.NewInt(int64(num))
	}
	return val
}

func InitLogLevel(cfg Config) error {
	if cfg.Log.Level != "" {
		lvl, err := log.ParseLevel(cfg.Log.Level)
		if err != nil {
			return err
		}
		log.SetLevel(lvl)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return nil
}

func GetConfigFromEnv() (Config, error) {
	cfgJson := os.Getenv(EnvConfig)
	if cfgJson == "" {
		return Config{}, fmt.Errorf("%s is required", EnvConfig)
	}
	var cfg Config
	if err := json.Unmarshal([]byte(cfgJson), &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func GetConfig(filename string) (Config, error) {
	var cfg Config
	if err := readFile(filename, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func readFile(filename string, cfg *Config) error {
	f, err := os.Open(filename)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(cfg)
}

// EnvDefaults contain default values for one env.
type EnvDefaults struct {
	DiscoSubdomain string
}

// GetEnvDefaults returns the default values for an env.
func GetEnvDefaults(development bool) EnvDefaults {
	if development {
		return EnvDefaults{
			DiscoSubdomain: "disco-dev",
		}
	}
	return EnvDefaults{
		DiscoSubdomain: "disco",
	}
}

// ENS contains the default names.
type ENS struct {
	Dispatch string
	Alerts   string
	Agents   string
}

// GetENSNames returns the default ENS names.
func GetENSNames() *ENS {
	return &ENS{
		Dispatch: "dispatch.forta.eth",
		Alerts:   "alerts.forta.eth",
		Agents:   "agents.registries.forta.eth",
	}
}
