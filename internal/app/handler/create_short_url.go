package handler

import (
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func ManagerCreateShortURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		sugar.Errorln("CreateShortURL body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	UserID := req.Context().Value("userID").(int)
	urlMask := utils.CreateRandomMask()
	getService := shorten.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	urlMask, err = getService.Create(urlMask, string(body), UserID)
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
