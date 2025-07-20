package handler

import (
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	myProto "github.com/VladimirSh98/urlShortener/proto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strings"
)

// ManagerCreateShortURL create short URL by text request
func (h *Handler) ManagerCreateShortURL(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		sugar.Errorln("CreateShortURL body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userIDRaw := req.Context().Value(middleware.UserIDKey)
	UserID, ok := userIDRaw.(int)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlMask := utils.CreateRandomMask()
	urlMask, err = h.service.Create(urlMask, string(body), UserID)
	var resStatus int
	if errors.Is(err, customErr.ErrConstraintViolation) {
		resStatus = http.StatusConflict
	} else {
		resStatus = http.StatusCreated
	}
	responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, urlMask)
	accept := req.Header.Get("Accept")
	switch {
	case strings.Contains(accept, "application/grpc"):
		grpcResponse := &myProto.ShortenResponse{
			Result: responseURL,
		}
		var data []byte
		data, err = proto.Marshal(grpcResponse)
		if err != nil {
			sugar.Errorln("Failed to marshal gRPC response:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/grpc+proto")
		res.WriteHeader(resStatus)
		if _, err = res.Write(data); err != nil {
			sugar.Errorln("Failed to write gRPC response:", err)
		}

	default:
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(resStatus)
		if _, err = res.Write([]byte(responseURL)); err != nil {
			sugar.Errorln("CreateShortURL response error", err)
		}
	}
}
