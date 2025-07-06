package service

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/handler"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Run service
func Run() error {
	sugar := zap.S()
	sugar.Infoln("Prefill data success")
	repo := dbRepo.NewShortenRepository(database.DBConnection.Conn)
	service := shorten.NewShortenService(repo)
	customHandler := handler.NewHandler(service)
	err := Prefill(service)
	if err != nil {
		sugar.Warnln("Prefill data error %v", err)
	}
	router := chi.NewMux()
	router.Use(middleware.Config)
	router.Mount("/debug", chiMiddleware.Profiler())
	router.Get("/ping", customHandler.Ping)
	router.Post("/", customHandler.ManagerCreateShortURL)
	router.Post("/api/shorten", customHandler.ManagerCreateShortURLByJSON)
	router.Post("/api/shorten/batch", customHandler.ManagerCreateShortURLBatch)
	router.Get("/{id}", customHandler.ManagerReturnFullURL)
	router.Get("/api/user/urls", customHandler.ManagerGetURLsByUser)
	router.Delete("/api/user/urls", customHandler.ManagerDeleteURLsByID)
	if config.EnableHTTPS {
		return http.ListenAndServeTLS(config.FlagRunAddr, config.CertFile, config.KeyFile, router)
	} else {
		return http.ListenAndServe(config.FlagRunAddr, router)
	}
}
