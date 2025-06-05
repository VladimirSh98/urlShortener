package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"go.uber.org/zap"
)

func (h *Handler) ManagerDeleteURLsByID(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		sugar.Errorln("ManagerDeleteURLsByID body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	UserID := req.Context().Value(middleware.UserIDKey).(int)
	var data []string
	err = json.Unmarshal(body, &data)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON json unmarshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	go h.service.BatchUpdate(data, UserID)
	res.WriteHeader(http.StatusAccepted)
}
