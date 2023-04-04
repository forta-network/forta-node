package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/creasty/defaults"
	"github.com/forta-network/forta-core-go/protocol/settings"
)

type PublicAPIProxyConfig struct {
	Url             string            `yaml:"url" json:"url" validate:"omitempty,url" default:"https://api.forta.network"`
	Headers         map[string]string `yaml:"headers" json:"headers"`
	RateLimitConfig *RateLimitConfig  `yaml:"rateLimit" json:"rateLimit"`
}
type JsonRpcConfig struct {
	Url     string            `yaml:"url" json:"url" validate:"omitempty,url"`
	Headers map[string]string `yaml:"headers" json:"headers"`
}

type ScannerConfig struct {
	JsonRpc              JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
	DisableAutostart     bool          `yaml:"disableAutostart" json:"disableAutostart"`
	BlockRateLimit       int           `yaml:"blockRateLimit" json:"blockRateLimit" default:"200"`
	BlockMaxAgeSeconds   int64         `yaml:"blockMaxAgeSeconds" json:"blockMaxAgeSeconds" default:"600"`
	RetryIntervalSeconds int64         `yaml:"retryIntervalSeconds" json:"retryIntervalSeconds" default:"8"`
	AlertAPIURL          string        `yaml:"apiUrl" json:"apiUrl" default:"https://api.forta.network/graphql" validate:"url"`
}

type TraceConfig struct {
	JsonRpc JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
	Enabled bool          `yaml:"enabled" json:"enabled"`
}

type RateLimitConfig struct {
	Rate  float64 `yaml:"rate" json:"rate"`
	Burst int     `yaml:"burst" json:"burst" validate:"min=1"`
}

type JsonRpcProxyConfig struct {
	JsonRpc         JsonRpcConfig    `yaml:"jsonRpc" json:"jsonRpc"`
	RateLimitConfig *RateLimitConfig `yaml:"rateLimit" json:"rateLimit"`
}

type LogConfig struct {
	Level       string `yaml:"level" json:"level" default:"info" `
	MaxLogSize  string `yaml:"maxLogSize" json:"maxLogSize" default:"50m" `
	MaxLogFiles int    `yaml:"maxLogFiles" json:"maxLogFiles" default:"10" `
}

type RegistryConfig struct {
	ChainID              uint64        `yaml:"chainId" json:"chainId" default:"137"`
	JsonRpc              JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc" default:"{\"url\": \"https://rpc.ankr.com/polygon\"}"`
	IPFS                 IPFSConfig    `yaml:"ipfs" json:"ipfs"`
	ContainerRegistry    string        `yaml:"containerRegistry" json:"containerRegistry" validate:"hostname|hostname_port" default:"disco.forta.network" `
	Username             string        `yaml:"username" json:"username"`
	Password             string        `yaml:"password" json:"password"`
	Disable              bool          `yaml:"disable" json:"disable"` // for testing situations
	CheckIntervalSeconds int           `yaml:"checkIntervalSeconds" json:"checkIntervalSeconds" default:"15"`
}

type IPFSConfig struct {
	GatewayURL string `yaml:"gatewayUrl" json:"gatewayUrl" validate:"url" default:"https://ipfs.forta.network" `
	APIURL     string `yaml:"apiUrl" json:"apiUrl" validate:"url" default:"https://ipfs.forta.network" `
	Username   string `yaml:"username" json:"username"`
	Password   string `yaml:"password" json:"password"`
}

type BatchConfig struct {
	SkipEmpty                    bool `yaml:"skipEmpty" json:"skipEmpty"`
	IntervalSeconds              *int `yaml:"intervalSeconds" json:"intervalSeconds" default:"15"`
	MetricsBucketIntervalSeconds *int `yaml:"metricsBucketIntervalSeconds" json:"metricsBucketIntervalSeconds" default:"60"`
	MaxAlerts                    *int `yaml:"maxAlerts" json:"maxAlerts" default:"1000" `
}

type PublisherConfig struct {
	SkipPublish   bool        `yaml:"skipPublish" json:"skipPublish" default:"false"`
	AlwaysPublish bool        `yaml:"alwaysPublish" json:"alwaysPublish" default:"false"`
	APIURL        string      `yaml:"apiUrl" json:"apiUrl" default:"https://alerts.forta.network" validate:"url"`
	IPFS          IPFSConfig  `yaml:"ipfs" json:"ipfs" validate:"required_unless=SkipPublish true"`
	Batch         BatchConfig `yaml:"batch" json:"batch"`
}

type ResourcesConfig struct {
	DisableAgentLimits bool    `yaml:"disableAgentLimits" json:"disableAgentLimits" default:"false" `
	AgentMaxMemoryMiB  int     `yaml:"agentMaxMemoryMib" json:"agentMaxMemoryMib" validate:"omitempty,min=100"`
	AgentMaxCPUs       float64 `yaml:"agentMaxCpus" json:"agentMaxCpus" validate:"omitempty,gt=0"`
}

type ENSConfig struct {
	DefaultContract bool   `yaml:"defaultContract" json:"defaultContract" default:"false" `
	ContractAddress string `yaml:"contractAddress" json:"contractAddress" validate:"omitempty,eth_addr" default:"0x08f42fcc52a9C2F391bF507C4E8688D0b53e1bd7"`
	Override        bool   `yaml:"override" json:"override" default:"false"`
}

type TelemetryConfig struct {
	URL       string `yaml:"url" json:"url" default:"https://alerts.forta.network/telemetry" validate:"url"`
	CustomURL string `yaml:"customUrl" validate:"omitempty,url"`
	Disable   bool   `yaml:"disable" json:"disable"`
}

type AutoUpdateConfig struct {
	Disable              bool `yaml:"disable" json:"disable"`
	UpdateDelay          *int `yaml:"updateDelay" json:"updateDelay"`
	TrackPrereleases     bool `yaml:"trackPrereleases" json:"trackPrereleases"`
	CheckIntervalSeconds int  `yaml:"checkIntervalSeconds" json:"checkIntervalSeconds" default:"60"` // 1m
}

type AgentLogsConfig struct {
	URL     string `yaml:"url" json:"url" default:"https://alerts.forta.network/logs/agents" validate:"url"`
	Disable bool   `yaml:"disable" json:"disable"`
}

type ContainerRegistryConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type RuntimeLimits struct {
	StartBlock         *uint64 `yaml:"startBlock" json:"startBlock"`
	StopBlock          *uint64 `yaml:"stopBlock" json:"stopBlock" validate:"omitempty,gtfield=StartBlock"`
	StopTimeoutSeconds int     `yaml:"stopTimeoutSeconds" json:"stopTimeoutSeconds" default:"30"`
	StartCombiner      uint64  `yaml:"startCombiner" json:"startCombiner"`
	StopCombiner       uint64  `yaml:"stopCombiner" json:"stopCombiner"`
}

type RedisConfig struct {
	Address  string `yaml:"address" json:"address"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

type RedisClusterConfig struct {
	Addresses []string `yaml:"addresses" json:"addresses"`
	Password  string   `yaml:"password" json:"password"`
	DB        int      `yaml:"db" json:"db"`
}

type DeduplicationConfig struct {
	TTLSeconds   int                 `yaml:"ttlSeconds" json:"ttlSeconds" default:"300"`
	Redis        *RedisConfig        `yaml:"redis" json:"redis"`
	RedisCluster *RedisClusterConfig `yaml:"redisCluster" json:"redisCluster"`
}

type StandaloneModeConfig struct {
	Enable        bool     `yaml:"enable" json:"enable"`
	BotContainers []string `yaml:"botContainers" json:"botContainers"`
}

type LocalModeConfig struct {
	Enable                bool                     `yaml:"enable" json:"enable"`
	IncludeMetrics        bool                     `yaml:"includeMetrics" json:"includeMetrics"`
	BotIDs                []string                 `yaml:"botIds" json:"botIds"`
	BotImages             []string                 `yaml:"botImages" json:"botImages"`
	WebhookURL            string                   `yaml:"webhookUrl" json:"webhookUrl"`
	LogFileName           string                   `yaml:"logFileName" json:"logFileName"`
	LogToStdout           bool                     `yaml:"logToStdout" json:"logToStdout"`
	ContainerRegistry     *ContainerRegistryConfig `yaml:"containerRegistry" json:"containerRegistry"`
	RuntimeLimits         RuntimeLimits            `yaml:"runtimeLimits" json:"runtimeLimits"`
	ForceEnableInspection bool                     `yaml:"forceEnableInspection" json:"forceEnableInspection"`
	Deduplication         *DeduplicationConfig     `yaml:"deduplication" json:"deduplication"`
	ShardedBots           []*LocalShardedBot       `yaml:"shardedBots" json:"shardedBots"`
	PrivateKeyHex         string                   `yaml:"privateKeyHex" json:"privateKeyHex"`
	Standalone            StandaloneModeConfig     `yaml:"standalone" json:"standalone"`
}

// IsStandalone checks if the node is in standalone mode. It should only be available
// as another local mode setting.
func (lmc LocalModeConfig) IsStandalone() bool {
	return lmc.Enable && lmc.Standalone.Enable
}

type LocalShardedBot struct {
	BotImage *string `yaml:"botImage" json:"botImage"`
	// number of shards for bot
	Shards uint `yaml:"shards" json:"shards"`
	// target per shard for bot
	Target uint `yaml:"target" json:"target"`
}

type InspectionConfig struct {
	BlockInterval     *int `yaml:"blockInterval" json:"blockInterval"`
	NetworkSavingMode bool `yaml:"networkSavingMode" json:"networkSavingMode"`
	InspectAtStartup  bool `yaml:"inspectAtStartup" json:"inspectAtStartup" default:"true"`
}

type StorageConfig struct {
	Provide string `yaml:"provide" json:"provide" default:"https://ipfs-router.forta.network/provide"`
	Reframe string `yaml:"reframe" json:"reframe" default:"https://ipfs-router.forta.network/reframe"`
}

type CombinerConfig struct {
	AlertAPIURL       string `yaml:"alertApiUrl" json:"alertApiUrl" default:"http://forta-public-api:8535" validate:"url"`
	CombinerCachePath string `yaml:"alertCachePath" json:"alert_cache_path"`
}

type AdvancedConfig struct {
	SafeOffset bool `yaml:"safeOffset" json:"safeOffset"`
}

type Config struct {
	// runtime values

	Development bool   `yaml:"-" json:"_development"`
	FortaDir    string `yaml:"-" json:"_fortaDir"`
	KeyDirPath  string `yaml:"-" json:"_keyDirPath"`
	Passphrase  string `yaml:"-" json:"_passphrase"`

	// yaml config values

	ChainID int `yaml:"chainId" json:"chainId" default:"1" `

	Scan  ScannerConfig `yaml:"scan" json:"scan"`
	Trace TraceConfig   `yaml:"trace" json:"trace"`

	Registry         RegistryConfig       `yaml:"registry" json:"registry"`
	Publish          PublisherConfig      `yaml:"publish" json:"publish"`
	JsonRpcProxy     JsonRpcProxyConfig   `yaml:"jsonRpcProxy" json:"jsonRpcProxy"`
	PublicAPIProxy   PublicAPIProxyConfig `yaml:"publicApiProxy" json:"publicApiProxy"`
	Log              LogConfig            `yaml:"log" json:"log"`
	ResourcesConfig  ResourcesConfig      `yaml:"resources" json:"resources"`
	ENSConfig        ENSConfig            `yaml:"ens" json:"ens"`
	TelemetryConfig  TelemetryConfig      `yaml:"telemetry" json:"telemetry"`
	AutoUpdate       AutoUpdateConfig     `yaml:"autoUpdate" json:"autoUpdate"`
	AgentLogsConfig  AgentLogsConfig      `yaml:"agentLogs" json:"agentLogs"`
	LocalModeConfig  LocalModeConfig      `yaml:"localMode" json:"localMode"`
	InspectionConfig InspectionConfig     `yaml:"inspection" json:"inspection"`
	StorageConfig    StorageConfig        `yaml:"storage" json:"storage"`
	CombinerConfig   CombinerConfig       `yaml:"combiner" json:"combiner"`
	AdvancedConfig   AdvancedConfig       `yaml:"advanced" json:"advanced"`
}

func (cfg *Config) ConfigFilePath() string {
	return path.Join(cfg.FortaDir, DefaultConfigFileName)
}

// GetConfigForContainer is how a container gets the forta configuration (file or env var)
func GetConfigForContainer() (Config, error) {
	cfg, err := getConfigFromFile()
	if err != nil {
		return Config{}, err
	}
	applyContextDefaults(&cfg)

	// initialize combiner cache dump path if cache is persistent
	if cfg.CombinerConfig.CombinerCachePath != "" {
		_, err = os.Stat(cfg.CombinerConfig.CombinerCachePath)
		if !os.IsNotExist(err) {
			return cfg, err
		}

		if os.IsNotExist(err) {
			if err := os.WriteFile(cfg.CombinerConfig.CombinerCachePath, []byte("{}"), 0666); err != nil {
				return cfg, err
			}
		}
	}

	return cfg, nil
}

// apply defaults that apply in certain contexts
func applyContextDefaults(cfg *Config) {
	chainSettings := settings.GetChainSettings(cfg.ChainID)
	if chainSettings.EnableTrace {
		cfg.Trace.Enabled = true
	}
	if cfg.ENSConfig.DefaultContract {
		cfg.ENSConfig.ContractAddress = ""
	}
	cfg.FortaDir = DefaultContainerFortaDirPath
	cfg.KeyDirPath = path.Join(cfg.FortaDir, DefaultKeysDirName)
	cfg.CombinerConfig.CombinerCachePath = path.Join(cfg.FortaDir, DefaultCombinerCacheFileName)
}

func getConfigFromFile() (cfg Config, err error) {
	var (
		successfullyLoadedTimes int
		wrappedErr              error
	)

	// if the default config file exists, load from there
	if err = checkIfConfigFileExists(DefaultContainerConfigPath); err == nil {
		if err = readYamlFile(DefaultContainerConfigPath, &cfg); err != nil {
			return
		}
		successfullyLoadedTimes++
	}
	if err != nil {
		wrappedErr = err
	}

	// if the wrapped config file exists, load from there
	if err = checkIfConfigFileExists(DefaultContainerWrappedConfigPath); err == nil {
		var wrapped map[string]Config
		if err = readYamlFile(DefaultContainerWrappedConfigPath, &wrapped); err != nil {
			return
		}
		var found bool
		cfg, found = wrapped[DefaultConfigWrapperKey]
		if !found {
			err = fmt.Errorf("wrapped config file was found but did not have the config under '%s'", DefaultConfigWrapperKey)
			return
		}
		successfullyLoadedTimes++
	}
	if err != nil {
		wrappedErr = fmt.Errorf("%v, %v", wrappedErr, err)
	}

	// at this point we expect that at least one of the config files were loaded without errors
	switch successfullyLoadedTimes {
	case 2:
		err = errors.New("multiple config files found in the forta dir - please use only one of them")
		return

	case 0:
		err = fmt.Errorf("failed to load any of the config files: %v", wrappedErr)
		return

	case 1:
		// yay! (ignore)

	default:
		err = fmt.Errorf("successfully loaded unexpected amount of config files (%d) - errors: %w", successfullyLoadedTimes, wrappedErr)
	}

	// finally set the defaults
	err = defaults.Set(&cfg)
	return
}

func checkIfConfigFileExists(configPath string) error {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return errors.New("config file not found")
	}
	if err != nil {
		return fmt.Errorf("failed to check if config file exists: %w", err)
	}
	return nil
}
