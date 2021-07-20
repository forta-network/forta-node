package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"OpenZeppelin/fortify-node/utils"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const EnvFortifyConfig = "FORTIFY_CONFIG"
const EnvJsonRpcHost = "JSON_RPC_HOST"
const EnvJsonRpcPort = "JSON_RPC_PORT"
const EnvAgentGrpcPort = "AGENT_GRPC_PORT"
const FortifyPrefix = "fortify"

const (
	//TODO: figure out protocol for this other than direct communication
	EnvQueryNode = "QUERY_NODE"
	// NATS node name env var
	EnvNatsHost = "NATS_HOST"
)

// Docker container names
var (
	DockerNatsContainerName         = fmt.Sprintf("%s-nats", FortifyPrefix)
	DockerScannerContainerName      = fmt.Sprintf("%s-scanner", FortifyPrefix)
	DockerJSONRPCProxyContainerName = fmt.Sprintf("%s-json-rpc", FortifyPrefix)
	DockerQueryContainerName        = fmt.Sprintf("%s-query", FortifyPrefix)

	DockerNetworkName = DockerScannerContainerName
)

// Global constant values
var (
	DefaultNatsPort    = "4222"
	DefaultIPFSGateway = "https://cloudflare-ipfs.com"
)

type AgentConfig struct {
	ID         string  `yaml:"id" json:"id"`
	Image      string  `yaml:"image" json:"image"`
	StartBlock *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock  *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
}

func (ac AgentConfig) ImageHash() string {
	_, digest := utils.SplitImageRef(ac.Image)
	return digest
}

func (ac AgentConfig) ContainerName() string {
	_, digest := utils.SplitImageRef(ac.Image)
	return fmt.Sprintf("fortify-agent-%s-%s", utils.ShortenString(ac.ID, 8), utils.ShortenString(digest, 4))
}

func (ac AgentConfig) GrpcPort() string {
	return "50051"
}

type DBConfig struct {
	Path string `yaml:"path" json:"path"`
}

type EthereumConfig struct {
	JsonRpcUrl string            `yaml:"jsonRpcUrl" json:"jsonRpcUrl"`
	Headers    map[string]string `yaml:"headers" json:"headers"`
}

type QueryConfig struct {
	QueryImage string   `yaml:"queryImage" json:"queryImage"`
	Port       int      `yaml:"port" json:"port"`
	DB         DBConfig `yaml:"db" json:"db"`
}

type ScannerConfig struct {
	ChainID      int            `yaml:"chainId" json:"chainId"`
	ScannerImage string         `yaml:"scannerImage" json:"scannerImage"`
	StartBlock   int            `yaml:"startBlock" json:"startBlock"`
	EndBlock     int            `yaml:"endBlock" json:"endBlock"`
	Ethereum     EthereumConfig `yaml:"ethereum" json:"ethereum"`
}

type TraceConfig struct {
	Ethereum EthereumConfig `yaml:"ethereum" json:"ethereum"`
	Enabled  bool           `yaml:"enabled" json:"enabled"`
}

type JsonRpcProxyConfig struct {
	JsonRpcImage string         `yaml:"jsonRpcImage" json:"jsonRpcImage"`
	Ethereum     EthereumConfig `yaml:"ethereum" json:"ethereum"`
}

type LogConfig struct {
	Level       string `yaml:"level" json:"level"`
	MaxLogSize  string `yaml:"maxLogSize" json:"maxLogSize"`
	MaxLogFiles int    `yaml:"maxLogFiles" json:"maxLogFiles"`
}

type RegistryConfig struct {
	Ethereum          EthereumConfig `yaml:"ethereum" json:"ethereum"`
	IPFSGateway       *string        `yaml:"ipfsGateway" json:"ipfs,omitempty"`
	ContractAddress   string         `yaml:"contractAddress" json:"contractAddress"`
	ContainerRegistry string         `yaml:"containerRegistry" json:"containerRegistry"`
	PoolID            string         `yaml:"poolId" json:"poolId"`
}

type Config struct {
	Registry     RegistryConfig     `yaml:"registry" json:"registry"`
	Scanner      ScannerConfig      `yaml:"scanner" json:"scanner"`
	Query        QueryConfig        `yaml:"query" json:"query"`
	Trace        TraceConfig        `yaml:"trace" json:"trace"`
	JsonRpcProxy JsonRpcProxyConfig `yaml:"json-rpc-proxy" json:"jsonRpcProxy"`
	Log          LogConfig          `yaml:"log" json:"log"`
}

func GetFortifyCfgDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.fortify", home), nil
}

func GetKeyStorePath() (string, error) {
	cfgDir, err := GetFortifyCfgDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.keys", cfgDir), nil
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
	cfgJson := os.Getenv(EnvFortifyConfig)
	if cfgJson == "" {
		return Config{}, fmt.Errorf("%s is required", EnvFortifyConfig)
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
