package handler

import (
	"encoding/json"
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

func ManagerDeleteURLsByID(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		sugar.Errorln("ManagerDeleteURLsByID body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var cookie *http.Cookie
	cookie, err = req.Cookie("userID")
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var UserID int
	UserID, err = strconv.Atoi(cookie.Value)
	if err != nil {
		sugar.Errorln("CreateShortURL convert cookie error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var data []string
	err = json.Unmarshal(body, &data)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON json unmarshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	go repository.BatchUpdate(data, UserID)
	res.WriteHeader(http.StatusAccepted)
}
