package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
	"os"
)

const ShortURLLength = 8

var (
	FlagResultAddr string
	FlagRunAddr    string
	DBFilePath     string
	DBFileName     string
)

func ParseFlags() error {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	defaultConfigValues, err := parseDefaultConfigValues()
	if err != nil {
		return err
	}
	flag.StringVar(&FlagRunAddr, "a", defaultConfigValues.ServerAddress, "Run address")
	flag.StringVar(&FlagResultAddr, "b", defaultConfigValues.BaseURL, "Result address")
	flag.StringVar(&DBFilePath, "f", defaultConfigValues.DBFilePath, "DB file path")
	flag.Parse()
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
		DBFilePath = defaultConfigValues.DBFilePath
	}
	DBFileName = defaultConfigValues.DBFileName
	return nil
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
