package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const EnvLogLevel = "LOG_LEVEL"

type AgentConfig struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

type Config struct {
	Zephyr struct {
		NodeImage string `yaml:"nodeImage"`
	} `yaml:"zephyr"`
	Ethereum struct {
		JsonRpcUrl string `yaml:"jsonRpcUrl"`
		StartBlock int    `yaml:"startBlock"`
	} `yaml:"ethereum"`
	Agents []AgentConfig `yaml:"agents"`
	Log    struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
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
