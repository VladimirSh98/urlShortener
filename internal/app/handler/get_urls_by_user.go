package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"go.uber.org/zap"
)

// ManagerGetURLsByUser get all URLs by user ID
func (h *Handler) ManagerGetURLsByUser(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	userIDRaw := req.Context().Value(middleware.UserIDKey)
	UserID, ok := userIDRaw.(int)
	if !ok {
		sugar.Errorln("invalid or missing user ID")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
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
	var response []getByUserIDResponseAPI
	for _, result := range results {
		ShortURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, result.ID)
		response = append(
			response,
			getByUserIDResponseAPI{ShortURL: ShortURL, URL: result.OriginalURL},
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
