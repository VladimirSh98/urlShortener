package service

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/handler"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Run() error {
	router := chi.NewMux()
	router.Post("/", handler.CreateShortURL)
	router.Get("/{id}", handler.ReturnFullURL)

	return http.ListenAndServe(config.FlagRunAddr, router)
}
