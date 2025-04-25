package handler

import (
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten_service"
	"go.uber.org/zap"
	"net/http"
)

func ManagerReturnFullURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	urlID := req.PathValue("id")
	resultURL, ok := memory.Get(urlID)
	if !ok {
		sugar.Infoln("ReturnFullURL no data by urlId", urlID)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	getService := shorten_service.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	recordFromDB, err := getService.GetByShortURL(urlID)
	if err != nil {
		sugar.Warnln("ReturnFullURL database error", err)
	}
	if recordFromDB.Archived {
		res.WriteHeader(http.StatusGone)
		return
	}
	res.Header().Set("Location", resultURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
