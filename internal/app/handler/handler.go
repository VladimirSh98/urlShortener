package handler

import (
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"net/http"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	err := database.DBConnection.Ping()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}

func CreateShortURL(res http.ResponseWriter, req *http.Request) {
	ManagerCreateShortURL(res, req)
}

func ReturnFullURL(res http.ResponseWriter, req *http.Request) {
	ManagerReturnFullURL(res, req)
}

func CreateShortURLByJSON(res http.ResponseWriter, req *http.Request) {
	ManagerCreateShortURLByJSON(res, req)
}

func CreateShortURLBatch(res http.ResponseWriter, req *http.Request) {
	ManagerCreateShortURLBatch(res, req)
}

func GetURLsByUser(res http.ResponseWriter, req *http.Request) {
	ManagerGetURLsByUser(res, req)
}

func DeleteURLsByID(res http.ResponseWriter, req *http.Request) {
	ManagerDeleteURLsByID(res, req)
}
