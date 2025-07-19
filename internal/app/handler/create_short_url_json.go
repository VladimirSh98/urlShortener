package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strings"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	myProto "github.com/VladimirSh98/urlShortener/proto"
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
	responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, urlMask)
	accept := req.Header.Get("Accept")

	switch {
	case strings.Contains(accept, "application/grpc"):
		grpcResponse := &myProto.ShortenResponse{
			Result: responseURL,
		}
		var grpcData []byte
		grpcData, err = proto.Marshal(grpcResponse)
		if err != nil {
			sugar.Errorln("Failed to marshal gRPC response:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/grpc+proto")
		res.WriteHeader(http.StatusOK)
		if _, err = res.Write(grpcData); err != nil {
			sugar.Errorln("Failed to write gRPC response:", err)
		}

	default:
		var response []byte
		response, err = json.Marshal(shortenResponseDataAPI{Result: responseURL})
		if err != nil {
			sugar.Errorln("JSON marshal error:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		if errors.Is(err, customErr.ErrConstraintViolation) {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}
		if _, err = res.Write(response); err != nil {
			sugar.Errorln("Failed to write JSON response:", err)
		}
	}
}
