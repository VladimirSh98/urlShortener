package main

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/logger"
	"github.com/VladimirSh98/urlShortener/internal/app/service"
	"log"
)

func main() {
	initLogger, err := logger.Initialize()
	defer initLogger.Sync()
	if err != nil {
		log.Fatalf("Logger configuration failed: %v", err)
	}
	err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Server configuration failed: %v", err)
	}
	if err = database.OpenConnectionDB(); err != nil {
		log.Printf("Database connection failed: %v", err)
	}
	err = service.Run()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
