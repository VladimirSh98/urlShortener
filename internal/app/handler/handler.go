package handler

import (
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"io"
	"net/http"
)

func CreateShortURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlMask := utils.CreateRandomMask()
	repository.Create(urlMask, string(body))
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, urlMask)
	_, err = res.Write([]byte(responseURL))
	if err != nil {
		return
	}
}

func ReturnFullURL(res http.ResponseWriter, req *http.Request) {
	urlID := req.PathValue("id")
	resultURL, ok := repository.Get(urlID)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", resultURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
