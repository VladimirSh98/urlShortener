package handler

import (
	"encoding/json"
	"fmt"
	myProto "github.com/VladimirSh98/urlShortener/proto"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"go.uber.org/zap"
)

// ManagerGetURLsByUser get all URLs by user ID
func (h *Handler) ManagerGetURLsByUser(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	userIDRaw := req.Context().Value(middleware.UserIDKey)
	UserID, ok := userIDRaw.(int)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := h.service.GetByUserID(UserID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	if len(results) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	var response []getByUserIDResponseAPI
	for _, result := range results {
		ShortURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, result.ID)
		response = append(
			response,
			getByUserIDResponseAPI{ShortURL: ShortURL, URL: result.OriginalURL},
		)
	}

	accept := req.Header.Get("Accept")

	switch {
	case strings.Contains(accept, "application/grpc"):
		h.handleGRPCGetURLsByUser(res, response, sugar)
	default:
		h.handleJSONGetURLsByUser(res, response, sugar)
	}
}

func (h *Handler) handleJSONGetURLsByUser(res http.ResponseWriter, response []getByUserIDResponseAPI, sugar *zap.SugaredLogger) {
	if len(response) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		sugar.Warnln("ManagerGetURLsByUser json marshall error", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = res.Write(jsonResponse); err != nil {
		sugar.Errorln("ManagerGetURLsByUser response error", err)
		res.WriteHeader(http.StatusBadRequest)
	}
}

func (h *Handler) handleGRPCGetURLsByUser(
	res http.ResponseWriter,
	response []getByUserIDResponseAPI,
	sugar *zap.SugaredLogger,
) {
	grpcResponse := &myProto.GetUserURLsResponse{
		Urls: make([]*myProto.UserURLResponse, 0, len(response)),
	}

	for _, result := range response {
		grpcResponse.Urls = append(grpcResponse.Urls, &myProto.UserURLResponse{
			ShortUrl:    result.ShortURL,
			OriginalUrl: result.URL,
		})
	}
	data, err := proto.Marshal(grpcResponse)
	if err != nil {
		sugar.Errorln("Failed to marshal gRPC response:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/grpc+proto")
	res.WriteHeader(http.StatusOK)

	if _, err = res.Write(data); err != nil {
		sugar.Errorln("Failed to write gRPC response:", err)
	}
}
