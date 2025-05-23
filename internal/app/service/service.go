package service

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/handler"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func Run() error {
	err := Prefill()
	sugar := zap.S()
	if err != nil {
		sugar.Warnln("Prefill data error %v", err)
	}
	sugar.Infoln("Prefill data success")
	repo := dbRepo.ShortenRepository{Conn: database.DBConnection.Conn}
	service := shorten.NewShortenService(repo)
	customHandler := handler.NewHandler(service)
	router := chi.NewMux()
	router.Use(middleware.Config)
	router.Get("/ping", customHandler.Ping)
	router.Post("/", customHandler.ManagerCreateShortURL)
	router.Post("/api/shorten", customHandler.ManagerCreateShortURLByJSON)
	router.Post("/api/shorten/batch", customHandler.ManagerCreateShortURLBatch)
	router.Get("/{id}", customHandler.ManagerReturnFullURL)
	router.Get("/api/user/urls", customHandler.ManagerGetURLsByUser)
	router.Delete("/api/user/urls", customHandler.ManagerDeleteURLsByID)

	return http.ListenAndServe(config.FlagRunAddr, router)
}
