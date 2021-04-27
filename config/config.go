package config

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const EnvZephyrConfig = "ZEPHYR_CONFIG"
const EnvAgentKeys = "AGENT_KEYS"
const EnvAgentKey = "AGENT_KEY"
const EnvAgentProxy = "AGENT_PROXY"
const ZephyrPrefix = "zephyr"

type AgentConfig struct {
	Name  string `yaml:"name" json:"name"`
	Image string `yaml:"image" json:"image"`
}

func (ac AgentConfig) ContainerName() string {
	return fmt.Sprintf("%s-agent-%s", ZephyrPrefix, ac.Name)
}

type Config struct {
	Zephyr struct {
		NodeImage  string `yaml:"nodeImage" json:"nodeImage"`
		ProxyImage string `yaml:"proxyImage" json:"proxyImage"`
	} `yaml:"zephyr" json:"zephyr"`
	Ethereum struct {
		JsonRpcUrl string `yaml:"jsonRpcUrl" json:"jsonRpcUrl"`
		StartBlock int    `yaml:"startBlock" json:"startBlock"`
	} `yaml:"ethereum" json:"ethereum"`
	Agents []AgentConfig `yaml:"agents" json:"agents"`
	Log    struct {
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
	cfgJson := os.Getenv(EnvZephyrConfig)
	if cfgJson == "" {
		return Config{}, fmt.Errorf("%s is required", EnvZephyrConfig)
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
