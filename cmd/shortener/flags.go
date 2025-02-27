package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

var flagRunAddr, flagResultAddr string

func parseFlags() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "Run address")
	flag.StringVar(&flagResultAddr, "b", "http://localhost:8080", "Result address")
	flag.Parse()
	if cfg.ServerAddress != "" {
		flagRunAddr = cfg.ServerAddress
	}
	if cfg.BaseURL != "" {
		flagResultAddr = cfg.BaseURL
	}
}
