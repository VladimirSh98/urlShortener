package handler

import (
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
	"go.uber.org/zap"
	"net/http"
)

func ManagerReturnFullURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	urlID := req.PathValue("id")
	resultURL, ok := repository.Get(urlID)
	if !ok {
		sugar.Infoln("ReturnFullURL no data by urlId", urlID)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	recordFromDB, err := repository.GetByShortURLFromBD(urlID)
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
