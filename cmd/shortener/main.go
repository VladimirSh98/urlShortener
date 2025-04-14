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
	err = database.DBConnection.OpenConnection()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.DBConnection.CloseConnection()
	err = database.DBConnection.UpgradeMigrations()
	if err != nil {
		log.Printf("Database migrations failed: %v", err)
	}
	err = service.Run()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
