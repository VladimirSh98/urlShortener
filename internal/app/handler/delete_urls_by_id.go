package handler

import (
	"encoding/json"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"

	"go.uber.org/zap"
	"io"
	"net/http"
)

func ManagerDeleteURLsByID(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		sugar.Errorln("ManagerDeleteURLsByID body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	UserID := req.Context().Value("userID").(int)
	var data []string
	err = json.Unmarshal(body, &data)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON json unmarshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	getService := shorten.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	go getService.BatchUpdate(data, UserID)
	res.WriteHeader(http.StatusAccepted)
}
