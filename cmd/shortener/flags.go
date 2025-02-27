package main

import "flag"

var flagRunAddr, flagResultAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "Run address")
	flag.StringVar(&flagResultAddr, "b", "http://localhost:8080", "Result address")
	flag.Parse()
}
