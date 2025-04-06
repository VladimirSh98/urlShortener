package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func ManagerCreateShortURLByJSON(res http.ResponseWriter, req *http.Request) {
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
	urlMask, err = repository.Create(urlMask, data.URL)
	res.Header().Set("Content-Type", "application/json")
	if errors.Is(err, customErr.ErrConstraintViolation) {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
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
