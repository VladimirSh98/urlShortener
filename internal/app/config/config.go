package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
	"os"
)

const ShortURLLength = 8
const CreateShortURLPath = "/"
const GetShortURLPath = "/{id}"

var (
	FlagResultAddr string
	FlagRunAddr    string
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

type defaultConfig struct {
	ServerAddress string `yaml:"server_address"`
	BaseURL       string `yaml:"base_url"`
}

func ParseFlags() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	defaultConfigValues := parseDefaultConfigValues()
	flag.StringVar(&FlagRunAddr, "a", defaultConfigValues.ServerAddress, "Run address")
	flag.StringVar(&FlagResultAddr, "b", defaultConfigValues.BaseURL, "Result address")
	flag.Parse()
	if cfg.ServerAddress != "" {
		FlagRunAddr = cfg.ServerAddress
	}
	if cfg.BaseURL != "" {
		FlagResultAddr = cfg.BaseURL
	}
}

func parseDefaultConfigValues() defaultConfig {
	defaultData, err := os.ReadFile("default_config.yaml")
	if err != nil {
		panic(err)
	}

	var defaultConfigValues defaultConfig
	err = yaml.Unmarshal(defaultData, &defaultConfigValues)
	if err != nil {
		panic(err)
	}
	return defaultConfigValues
}
