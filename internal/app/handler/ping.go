package handler

import (
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten_service"
	"net/http"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	getService := shorten_service.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	err := getService.Ping()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}
