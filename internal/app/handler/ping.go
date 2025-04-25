package handler

import (
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
	"net/http"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	getService := shortenService.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	err := getService.Ping()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}
