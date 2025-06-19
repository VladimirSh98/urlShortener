package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

// ShortURLLength contains short URL length
const ShortURLLength = 8

// FlagResultAddr contains result link URL
var FlagResultAddr string

// FlagRunAddr contains project URL
var FlagRunAddr string

// DBFilePath contains files path for data saving
var DBFilePath string

// DatabaseDSN contains DB dsn
var DatabaseDSN string

// DefaultConfigValues contains default credentials
var DefaultConfigValues defaultConfig

// LoadConfig loads the project configuration
func LoadConfig() error {
	var cfg config
	var err error

	err = env.Parse(&cfg)
	if err != nil {
		return err
	}
	DefaultConfigValues, err = parseDefaultConfigValues()
	if err != nil {
		return err
	}
	parseFlag()
	if cfg.ServerAddress != "" {
		FlagRunAddr = cfg.ServerAddress
	}
	if cfg.BaseURL != "" {
		FlagResultAddr = cfg.BaseURL
	}
	if cfg.DBFilePath != "" {
		DBFilePath = cfg.DBFilePath
	}
	if DBFilePath == "" {
		DBFilePath = DefaultConfigValues.DBFilePath
	}
	if cfg.DatabaseDSN != "" {
		DatabaseDSN = cfg.DatabaseDSN
	}
	return nil
}

func parseFlag() {
	flag.StringVar(&FlagRunAddr, "a", DefaultConfigValues.ServerAddress, "Run address")
	flag.StringVar(&FlagResultAddr, "b", DefaultConfigValues.BaseURL, "Result address")
	flag.StringVar(&DBFilePath, "f", DefaultConfigValues.DBFilePath, "DB file path")
	flag.StringVar(&DatabaseDSN, "d", DefaultConfigValues.DatabaseDSN, "DB path")
	flag.Parse()
}

func parseDefaultConfigValues() (defaultConfig, error) {
	defaultData, err := os.ReadFile("default_config.yaml")
	if err != nil {
		return defaultConfig{}, err
	}

	var defaultConfigValues defaultConfig
	err = yaml.Unmarshal(defaultData, &defaultConfigValues)
	if err != nil {
		return defaultConfig{}, err
	}
	return defaultConfigValues, nil
}
