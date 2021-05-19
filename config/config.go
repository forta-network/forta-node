package config

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const EnvFortifyConfig = "FORTIFY_CONFIG"
const EnvJsonRpcHost = "JSON_RPC_HOST"
const EnvJsonRpcPort = "JSON_RPC_PORT"
const FortifyPrefix = "fortify"

//TODO: figure out protocol for this other than direct communication
const EnvQueryNode = "QUERY_NODE"

type AgentConfig struct {
	Name      string `yaml:"name" json:"name"`
	Image     string `yaml:"image" json:"image"`
	ImageHash string `yaml:"imageHash" json:"imageHash"`
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
	ScannerImage string         `yaml:"scannerImage" json:"scannerImage"`
	StartBlock   int            `yaml:"startBlock" json:"startBlock"`
	Ethereum     EthereumConfig `yaml:"ethereum" json:"ethereum"`
}

type JsonRpcProxyConfig struct {
	JsonRpcImage string         `yaml:"jsonRpcImage" json:"jsonRpcImage"`
	Ethereum     EthereumConfig `yaml:"ethereum" json:"ethereum"`
}

func (ac AgentConfig) ContainerName() string {
	return fmt.Sprintf("%s-agent-%s", FortifyPrefix, ac.Name)
}

type Config struct {
	Scanner      ScannerConfig      `yaml:"scanner" json:"scanner"`
	Query        QueryConfig        `yaml:"query" json:"query"`
	JsonRpcProxy JsonRpcProxyConfig `yaml:"json-rpc-proxy" json:"jsonRpcProxy"`
	Agents       []AgentConfig      `yaml:"agents" json:"agents"`
	Log          struct {
		Level string `yaml:"level" json:"level"`
	} `yaml:"log" json:"log"`
}

func (c Config) AgentContainerNames() []string {
	var agents []string
	for _, agt := range c.Agents {
		agents = append(agents, agt.ContainerName())
	}
	return agents
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
