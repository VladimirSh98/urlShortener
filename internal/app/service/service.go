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
	err := prefill()
	sugar := zap.S()
	if err != nil {
		sugar.Warnln("Prefill data error %v", err)
	}
	sugar.Infoln("Prefill data success")
	router := chi.NewMux()
	router.Use(middleware.Config)
	router.Post("/", handler.CreateShortURL)
	router.Post("/api/shorten", handler.CreateShortURLByJSON)
	router.Get("/{id}", handler.ReturnFullURL)

	return http.ListenAndServe(config.FlagRunAddr, router)
}
