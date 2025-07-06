package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/logger"
	"github.com/VladimirSh98/urlShortener/internal/app/service"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	if os.Getenv("ENABLE_PPROF") == "true" {
		go func() {
			log.Println("Starting pprof server on :6060")
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	log.Printf("Build version: %s\n", buildVersion)
	log.Printf("Build date: %s\n", buildDate)
	log.Printf("Build commit: %s\n", buildCommit)
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
