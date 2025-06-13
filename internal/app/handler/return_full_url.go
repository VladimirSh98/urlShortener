package handler

import (
	"net/http"

	"github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
	"go.uber.org/zap"
)

// ManagerReturnFullURL return URL by short URL
func (h *Handler) ManagerReturnFullURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	urlID := req.PathValue("id")
	resultURL, ok := memory.Get(urlID)
	if !ok {
		sugar.Infoln("ReturnFullURL no data by urlId", urlID)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	recordFromDB, err := h.service.GetByShortURL(urlID)
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
