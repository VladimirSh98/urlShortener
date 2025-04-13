package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func ManagerGetURLsByUser(res http.ResponseWriter, req *http.Request) {
	var err error
	var cookie *http.Cookie

	sugar := zap.S()
	cookie, err = req.Cookie("userID")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var UserID int
	UserID, err = strconv.Atoi(cookie.Value)
	if err != nil {
		sugar.Errorln("ManagerGetURLsByUser convert cookie error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := repository.GetByUserID(UserID)
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
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = res.Write(jsonResponse)
	if err != nil {
		sugar.Errorln("ManagerGetURLsByUser response error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
