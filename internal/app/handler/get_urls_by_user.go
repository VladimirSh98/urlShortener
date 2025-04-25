package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
	"go.uber.org/zap"
	"net/http"
)

func ManagerGetURLsByUser(res http.ResponseWriter, req *http.Request) {
	var err error

	sugar := zap.S()
	UserID := req.Context().Value("userID").(int)
	getService := shorten.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	results, err := getService.GetByUserID(UserID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	if len(results) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	var response []APIGetByUserIDResponse
	for _, result := range results {
		ShortURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, result.ID)
		response = append(
			response,
			APIGetByUserIDResponse{ShortURL: ShortURL, URL: result.OriginalURL},
		)
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		sugar.Warnln("ManagerGetURLsByUser json marshall error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = res.Write(jsonResponse)
	if err != nil {
		sugar.Errorln("ManagerGetURLsByUser response error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
