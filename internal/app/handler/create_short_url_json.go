package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shortenService"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

func ManagerCreateShortURLByJSON(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var cookie *http.Cookie
	cookie, err = req.Cookie("userID")
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	var UserID int
	UserID, err = strconv.Atoi(cookie.Value)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON convert cookie error", err)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var data APIShortenRequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON json unmarshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	v := validator.New()
	err = v.Struct(data)
	if err != nil {
		sugar.Warnln("ManagerCreateShortURLByJSON validation error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlMask := utils.CreateRandomMask()
	getService := shortenService.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	urlMask, err = getService.Create(urlMask, data.URL, UserID)
	res.Header().Set("Content-Type", "application/json")
	if errors.Is(err, customErr.ErrConstraintViolation) {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, urlMask)
	response, err := json.Marshal(APIShortenResponseData{Result: responseURL})
	if err != nil {
		sugar.Warnln("ManagerCreateShortURLByJSON json marshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = res.Write(response)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON response error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}
