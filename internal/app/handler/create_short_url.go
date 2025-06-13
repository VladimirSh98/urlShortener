package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ManagerCreateShortURL create short URL by text request
func (h *Handler) ManagerCreateShortURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		sugar.Errorln("CreateShortURL body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	UserID := req.Context().Value(middleware.UserIDKey).(int)
	urlMask := utils.CreateRandomMask()
	urlMask, err = h.service.Create(urlMask, string(body), UserID)
	res.Header().Set("Content-Type", "text/plain")
	if errors.Is(err, customErr.ErrConstraintViolation) {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, urlMask)
	_, err = res.Write([]byte(responseURL))
	if err != nil {
		sugar.Errorln("CreateShortURL response error", err)
		return
	}
}
