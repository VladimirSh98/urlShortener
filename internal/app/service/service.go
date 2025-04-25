package service

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/handler"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
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
	router := chi.NewMux()
	router.Use(middleware.Config)
	router.Get("/ping", handler.Ping)
	router.Post("/", handler.ManagerCreateShortURL)
	router.Post("/api/shorten", handler.ManagerCreateShortURLByJSON)
	router.Post("/api/shorten/batch", handler.ManagerCreateShortURLBatch)
	router.Get("/{id}", handler.ManagerReturnFullURL)
	router.Get("/api/user/urls", handler.ManagerGetURLsByUser)
	router.Delete("/api/user/urls", handler.ManagerDeleteURLsByID)

	return http.ListenAndServe(config.FlagRunAddr, router)
}
