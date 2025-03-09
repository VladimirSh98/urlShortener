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
	if !(req.Method == http.MethodPost && req.URL.Path == "/") {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
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
	res.Write([]byte(responseURL))
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
