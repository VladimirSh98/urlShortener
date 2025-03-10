package main

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/service"
	"log"
)

func main() {
	err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Server configuration failed: %v", err)
	}
	err = service.Run()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
