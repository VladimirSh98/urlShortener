package handler

import (
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func CreateShortURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		sugar.Errorln("CreateShortURL body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlMask := utils.CreateRandomMask()
	urlMask, err = repository.Create(urlMask, string(body))
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

func ReturnFullURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	urlID := req.PathValue("id")
	resultURL, ok := repository.Get(urlID)
	if !ok {
		sugar.Infoln("ReturnFullURL no data by urlId", urlID)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", resultURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func CreateShortURLByJSON(res http.ResponseWriter, req *http.Request) {
	ManagerCreateShortURLByJSON(res, req)
}

func Ping(res http.ResponseWriter, req *http.Request) {
	err := database.DBConnection.Ping()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}

func CreateShortURLBatch(res http.ResponseWriter, req *http.Request) {
	ManagerCreateShortURLBatch(res, req)
}
