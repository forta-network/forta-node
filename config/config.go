package config

import (
	"errors"
	"os"
	"path"

	"github.com/creasty/defaults"
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
	Level       string `yaml:"level" json:"level" default:"info" `
	MaxLogSize  string `yaml:"maxLogSize" json:"maxLogSize" default:"50m" `
	MaxLogFiles int    `yaml:"maxLogFiles" json:"maxLogFiles" default:"10" `
}

type RegistryConfig struct {
	JsonRpc           JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc" default:"{\"url\": \"https://polygon-rpc.com\"}"`
	IPFS              IPFSConfig    `yaml:"ipfs" json:"ipfs"`
	ContractAddress   string        `yaml:"contractAddress" json:"contractAddress" validate:"eth_addr"`
	ContainerRegistry string        `yaml:"containerRegistry" json:"containerRegistry" validate:"hostname" default:"disco.forta.network" `
	Username          string        `yaml:"username" json:"username"`
	Password          string        `yaml:"password" json:"password"`
	Disabled          bool          `yaml:"disabled" json:"disabled"` // for testing situations
}

type IPFSConfig struct {
	GatewayURL string `yaml:"gatewayUrl" json:"gatewayUrl" validate:"url" default:"https://ipfs.forta.network" `
	APIURL     string `yaml:"apiUrl" json:"apiUrl" validate:"url" default:"https://ipfs.forta.network" `
	Username   string `yaml:"username" json:"username"`
	Password   string `yaml:"password" json:"password"`
}

type BatchConfig struct {
	SkipEmpty       bool `yaml:"skipEmpty" json:"skipEmpty"`
	IntervalSeconds *int `yaml:"intervalSeconds" json:"intervalSeconds" default:"15" `
	MaxAlerts       *int `yaml:"maxAlerts" json:"maxAlerts" default:"1000" `
}

type TestAlertsConfig struct {
	Disable    bool   `yaml:"disable" json:"disable"`
	WebhookURL string `yaml:"webhookUrl" json:"webhookUrl" validate:"omitempty,url"`
}

type PublisherConfig struct {
	SkipPublish     bool             `yaml:"skipPublish" json:"skipPublish" default:"false" `
	ContractAddress string           `yaml:"contractAddress" json:"contractAddress"`
	JsonRpc         JsonRpcConfig    `yaml:"jsonRpc" json:"jsonRpc" default:"{\"url\": \"https://polygon-rpc.com\"}" `
	IPFS            IPFSConfig       `yaml:"ipfs" json:"ipfs" validate:"required_unless=SkipPublish true"`
	Batch           BatchConfig      `yaml:"batch" json:"batch"`
	TestAlerts      TestAlertsConfig `yaml:"testAlerts" json:"testAlerts"`
	GasPriceGwei    int64            `yaml:"gasPriceGwei" json:"gasPriceGwei" default:"50" `
	GasLimit        uint64           `yaml:"gasLimit" json:"gasLimit" default:"50000" `
}

type ResourcesConfig struct {
	DisableAgentLimits bool    `yaml:"disableAgentLimits" json:"disableAgentLimits" default:"false" `
	AgentMaxMemoryMiB  int     `yaml:"agentMaxMemoryMib" json:"agentMaxMemoryMib" validate:"omitempty,min=100"`
	AgentMaxCPUs       float64 `yaml:"agentMaxCpus" json:"agentMaxCpus" validate:"omitempty,gt=0"`
}

type ENSConfig struct {
	DefaultContract bool          `yaml:"defaultContract" json:"defaultContract" default:"false" `
	ContractAddress string        `yaml:"contractAddress" json:"contractAddress" validate:"omitempty,eth_addr" default:"0x08f42fcc52a9C2F391bF507C4E8688D0b53e1bd7"`
	JsonRpc         JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc" default:"{\"url\": \"https://polygon-rpc.com\"}" `
}

type Config struct {
	ChainID                      int            `yaml:"chainId" json:"chainId" default:"1" `
	Development                  bool           `yaml:"-" json:"_development"`
	FortaDir                     string         `yaml:"-" json:"_fortaDir"`
	KeyDirPath                   string         `yaml:"-" json:"_keyDirPath"`
	Passphrase                   string         `yaml:"-" json:"_passphrase"`
	ExposeNats                   bool           `yaml:"-" json:"_exposeNats"`
	LocalAgentsPath              string         `yaml:"-" json:"_localAgentsPath"`
	LocalAgents                  []*AgentConfig `yaml:"-" json:"_localAgents"`
	AgentRegistryContractAddress string         `yaml:"-" json:"_agentRegistry"`

	Scan  ScannerConfig `yaml:"scan" json:"scan"`
	Trace TraceConfig   `yaml:"trace" json:"trace"`

	Registry        RegistryConfig      `yaml:"registry" json:"registry"`
	Publish         PublisherConfig     `yaml:"publish" json:"publish"`
	JsonRpcProxy    *JsonRpcProxyConfig `yaml:"jsonRpcProxy" json:"jsonRpcProxy"`
	Log             LogConfig           `yaml:"log" json:"log"`
	ResourcesConfig ResourcesConfig     `yaml:"resources" json:"resources"`
	ENSConfig       ENSConfig           `yaml:"ens" json:"ens"`
}

func (cfg *Config) ConfigFilePath() string {
	return path.Join(cfg.FortaDir, DefaultConfigFileName)
}

//GetConfigForContainer is how a container gets the forta configuration (file or env var)
func GetConfigForContainer() (Config, error) {
	var cfg Config
	if _, err := os.Stat(DefaultContainerConfigPath); os.IsNotExist(err) {
		return cfg, errors.New("config file not found")
	}
	cfg, err := getConfigFromFile(DefaultContainerConfigPath)
	if err != nil {
		return Config{}, err
	}
	applyContextDefaults(&cfg)
	return cfg, nil
}

// apply defaults that apply in certain contexts
func applyContextDefaults(cfg *Config) {
	if cfg.ChainID == 1 {
		cfg.Trace.Enabled = true
	}
	if cfg.ENSConfig.DefaultContract {
		cfg.ENSConfig.ContractAddress = ""
	}
	cfg.FortaDir = DefaultContainerFortaDirPath
	cfg.KeyDirPath = path.Join(cfg.FortaDir, DefaultKeysDirName)
}

func getConfigFromFile(filename string) (Config, error) {
	var cfg Config
	if err := readFile(filename, &cfg); err != nil {
		return Config{}, err
	}
	if err := defaults.Set(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
