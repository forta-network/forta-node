package config

import (
	"fmt"
	"os"

	"github.com/creasty/defaults"
	"github.com/goccy/go-json"
)

type JsonRpcConfig struct {
	Url     string            `yaml:"url" json:"url" validate:"omitempty,url"`
	Headers map[string]string `yaml:"headers" json:"headers"`
}

type ScannerConfig struct {
	StartBlock       int           `yaml:"startBlock" json:"startBlock"`
	EndBlock         int           `yaml:"endBlock" json:"endBlock"`
	JsonRpc          JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
	DisableAutostart bool          `yaml:"disableAutostart" json:"disableAutostart"`
	BlockRateLimit   int           `yaml:"blockRateLimit" json:"blockRateLimit"`
}

type TraceConfig struct {
	JsonRpc JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
	Enabled bool          `yaml:"enabled" json:"enabled"`
}

type JsonRpcProxyConfig struct {
	JsonRpc JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
}

type LogConfig struct {
	Level       string `default:"info" yaml:"level" json:"level"`
	MaxLogSize  string `default:"50m" yaml:"maxLogSize" json:"maxLogSize"`
	MaxLogFiles int    `default:"10" yaml:"maxLogFiles" json:"maxLogFiles"`
}

type RegistryConfig struct {
	JsonRpc           JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
	IPFS              IPFSConfig    `yaml:"ipfs" json:"ipfs"`
	ContractAddress   string        `yaml:"contractAddress" json:"contractAddress" validate:"eth_addr"`
	ContainerRegistry string        `default:"disco.forta.network" yaml:"containerRegistry" json:"containerRegistry" validate:"hostname"`
	Username          string        `yaml:"username" json:"username"`
	Password          string        `yaml:"password" json:"password"`
	Disabled          bool          `yaml:"disabled" json:"disabled"` // for testing situations
}

type IPFSConfig struct {
	GatewayURL string `default:"https://ipfs.forta.network" yaml:"gatewayUrl" json:"gatewayUrl" validate:"url"`
	APIURL     string `default:"https://ipfs.forta.network" yaml:"apiUrl" json:"apiUrl" validate:"url"`
	Username   string `yaml:"username" json:"username"`
	Password   string `yaml:"password" json:"password"`
}

type BatchConfig struct {
	SkipEmpty       bool `yaml:"skipEmpty" json:"skipEmpty"`
	IntervalSeconds *int `default:"15" yaml:"intervalSeconds" json:"intervalSeconds"`
	MaxAlerts       *int `default:"1000" yaml:"maxAlerts" json:"maxAlerts"`
}

type TestAlertsConfig struct {
	Disable    bool   `yaml:"disable" json:"disable"`
	WebhookURL string `yaml:"webhookUrl" json:"webhookUrl" validate:"omitempty,url"`
}

type PublisherConfig struct {
	SkipPublish     bool             `default:"false" yaml:"skipPublish" json:"skipPublish"`
	ContractAddress string           `yaml:"contractAddress" json:"contractAddress"`
	JsonRpc         JsonRpcConfig    `default:"{\"url\": \"https://polygon-rpc.com\"}" yaml:"jsonRpc" json:"jsonRpc"`
	IPFS            IPFSConfig       `yaml:"ipfs" json:"ipfs" validate:"required_unless=SkipPublish true"`
	Batch           BatchConfig      `yaml:"batch" json:"batch"`
	TestAlerts      TestAlertsConfig `yaml:"testAlerts" json:"testAlerts"`
	GasPriceGwei    int64            `default:"50" yaml:"gasPriceGwei" json:"gasPriceGwei"`
	GasLimit        uint64           `default:"50000" yaml:"gasLimit" json:"gasLimit"`
}

type ResourcesConfig struct {
	DisableAgentLimits bool    `default:"false" yaml:"disableAgentLimits" json:"disableAgentLimits"`
	AgentMaxMemoryMiB  int     `yaml:"agentMaxMemoryMib" json:"agentMaxMemoryMib" validate:"omitempty,min=100"`
	AgentMaxCPUs       float64 `yaml:"agentMaxCpus" json:"agentMaxCpus" validate:"omitempty,gt=0"`
}

type ENSConfig struct {
	ContractAddress string        `default:"0x08f42fcc52a9C2F391bF507C4E8688D0b53e1bd7" yaml:"contractAddress" json:"contractAddress" validate:"omitempty,eth_addr"`
	JsonRpc         JsonRpcConfig `default:"{\"url\": \"https://polygon-rpc.com\"}" yaml:"jsonRpc" json:"jsonRpc"`
}

type Config struct {
	ChainID                      int            `default:"1" yaml:"chainId" json:"chainId"`
	Development                  bool           `yaml:"-" json:"_development"`
	FortaDir                     string         `yaml:"-" json:"_fortaDir"`
	ConfigPath                   string         `yaml:"-" json:"_configPath"`
	KeyDirPath                   string         `yaml:"-" json:"_keyDirPath"`
	Passphrase                   string         `yaml:"-" json:"_passphrase"`
	ExposeNats                   bool           `yaml:"-" json:"_exposeNats"`
	LocalAgentsPath              string         `yaml:"-" json:"_localAgentsPath"`
	LocalAgents                  []*AgentConfig `yaml:"-" json:"localAgents"`
	AgentRegistryContractAddress string         `yaml:"-" json:"agentRegistry"`

	Scan  ScannerConfig `yaml:"scan" json:"scan"`
	Trace TraceConfig   `yaml:"trace" json:"trace"`

	Registry        RegistryConfig      `yaml:"registry" json:"registry"`
	Publish         PublisherConfig     `yaml:"publish" json:"publish"`
	JsonRpcProxy    *JsonRpcProxyConfig `yaml:"jsonRpcProxy" json:"jsonRpcProxy"`
	Log             LogConfig           `yaml:"log" json:"log"`
	ResourcesConfig ResourcesConfig     `yaml:"resources" json:"resources"`
	ENSConfig       ENSConfig           `yaml:"ens" json:"ens"`
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
	if err := defaults.Set(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func GetConfig(filename string) (Config, error) {
	var cfg Config
	if err := readFile(filename, &cfg); err != nil {
		return Config{}, err
	}
	if err := defaults.Set(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
