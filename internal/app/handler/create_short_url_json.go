package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ManagerCreateShortURLByJSON create short URL by JSON request
func (h *Handler) ManagerCreateShortURLByJSON(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		sugar.Errorln("ManagerCreateShortURLByJSON body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userIDRaw := req.Context().Value(middleware.UserIDKey)
	UserID, ok := userIDRaw.(int)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var data shortenRequestDataAPI
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
	urlMask, err = h.service.Create(urlMask, data.URL, UserID)
	res.Header().Set("Content-Type", "application/json")
	if errors.Is(err, customErr.ErrConstraintViolation) {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, urlMask)
	response, err := json.Marshal(shortenResponseDataAPI{Result: responseURL})
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
