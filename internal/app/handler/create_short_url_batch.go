package handler

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strings"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	myProto "github.com/VladimirSh98/urlShortener/proto"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ManagerCreateShortURLBatch batch create short URLs by json request
func (h *Handler) ManagerCreateShortURLBatch(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		sugar.Errorln("CreateShortURLBatch body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	userIDRaw := req.Context().Value(middleware.UserIDKey)
	UserID, ok := userIDRaw.(int)
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var data []shortenBatchRequestAPI
	err = json.Unmarshal(body, &data)
	if err != nil {
		sugar.Errorln("CreateShortURLBatch json unmarshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	v := validator.New()
	err = v.Var(data, "required,dive")
	if err != nil {
		sugar.Warnln("CreateShortURLBatch validation error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	var result []shortenBatchResponseAPI
	result = make([]shortenBatchResponseAPI, 0)
	dataWithMask := generateMaskForManyURLs(data)
	var prepareDataForBatch []dbRepo.ShortenBatchRequest
	prepareDataForBatch = make([]dbRepo.ShortenBatchRequest, 0)
	for _, record := range dataWithMask {
		prepareDataForBatch = append(prepareDataForBatch, dbRepo.ShortenBatchRequest{
			URL:    record.URL,
			Mask:   record.Mask,
			UserID: UserID,
		})
	}
	h.service.BatchCreate(prepareDataForBatch)
	accept := req.Header.Get("Accept")

	switch {
	case strings.Contains(accept, "application/grpc"):
		grpcResponse := &myProto.BatchShortenResponseList{
			Urls: make([]*myProto.BatchShortenResponse, 0, len(dataWithMask)),
		}
		for _, record := range dataWithMask {
			responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, record.Mask)
			grpcResponse.Urls = append(grpcResponse.Urls, &myProto.BatchShortenResponse{
				CorrelationId: record.CorrelationID,
				ShortUrl:      responseURL,
			})
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
		for _, record := range dataWithMask {
			responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, record.Mask)
			result = append(result, shortenBatchResponseAPI{
				CorrelationID: record.CorrelationID,
				URL:           responseURL,
			})
		}
		var response []byte
		response, err = json.Marshal(result)
		if err != nil {
			sugar.Errorln("JSON marshal error:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err = res.Write(response); err != nil {
			sugar.Errorln("Failed to write JSON response:", err)
		}
	}
}

func generateMaskForManyURLs(data []shortenBatchRequestAPI) []shortenBatchRequestWithMaskAPI {
	response := make([]shortenBatchRequestWithMaskAPI, 0)
	for _, record := range data {
		mask := utils.CreateRandomMask()
		response = append(response, shortenBatchRequestWithMaskAPI{
			shortenBatchRequestAPI: record,
			Mask:                   mask,
		})
	}
	return response
}
