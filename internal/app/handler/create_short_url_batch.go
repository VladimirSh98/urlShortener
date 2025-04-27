package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (h *Handler) ManagerCreateShortURLBatch(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		sugar.Errorln("CreateShortURLBatch body read error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	UserID := req.Context().Value(middleware.UserIDKey).(int)
	var data []APIShortenBatchRequest
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
	var result []APIShortenBatchResponse
	result = make([]APIShortenBatchResponse, 0)
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
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	for _, record := range dataWithMask {
		responseURL := fmt.Sprintf("%s/%s", config.FlagResultAddr, record.Mask)
		result = append(result, APIShortenBatchResponse{
			CorrelationID: record.CorrelationID,
			URL:           responseURL,
		})
	}

	response, err := json.Marshal(result)
	if err != nil {
		sugar.Warnln("CreateShortURLBatch json marshall error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = res.Write(response)
	if err != nil {
		sugar.Errorln("CreateShortURLBatch response error", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func generateMaskForManyURLs(data []APIShortenBatchRequest) []APIShortenBatchRequestWithMask {
	response := make([]APIShortenBatchRequestWithMask, 0)
	for _, record := range data {
		mask := utils.CreateRandomMask()
		response = append(response, APIShortenBatchRequestWithMask{
			APIShortenBatchRequest: record,
			Mask:                   mask,
		})
	}
	return response
}
