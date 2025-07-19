package handler

import (
	"encoding/json"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"go.uber.org/zap"
	"net/http"
)

// GetStats get stats for urls count and users count
func (h *Handler) GetStats(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	result, err := h.service.GetAllRecords()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := getStatsResponseAPI{
		URLS:  len(result),
		Users: int(middleware.UserCount),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		sugar.Warnln("GetStats json marshall error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = res.Write(jsonResponse)
	if err != nil {
		sugar.Errorln("GetStats response error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
