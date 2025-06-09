package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"go.uber.org/zap"
)

func (h *Handler) ManagerGetURLsByUser(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	UserID := req.Context().Value(middleware.UserIDKey).(int)
	results, err := h.service.GetByUserID(UserID)
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
