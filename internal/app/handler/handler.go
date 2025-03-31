package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/go-playground/validator/v10"
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
	repository.Create(urlMask, string(body), true)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
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
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		sugar.Errorln("CreateShortURLByJSON body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var data APIShortenRequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		sugar.Errorln("CreateShortURLByJSON json unmarshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	v := validator.New()
	err = v.Struct(data)
	if err != nil {
		sugar.Warnln("CreateShortURLByJSON validation error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlMask := utils.CreateRandomMask()
	repository.Create(urlMask, data.URL, true)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, urlMask)
	response, err := json.Marshal(APIShortenResponseData{Result: responseURL})
	if err != nil {
		sugar.Warnln("CreateShortURLByJSON json marshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = res.Write(response)
	if err != nil {
		sugar.Errorln("CreateShortURLByJSON response error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func Ping(res http.ResponseWriter, req *http.Request) {
	err := database.DBConnection.Ping()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}
