package handler

import "github.com/VladimirSh98/urlShortener/internal/app/service/shorten"

type Handler struct {
	service shorten.ShortenServiceInterface
}

func NewHandler(service shorten.ShortenServiceInterface) *Handler {
	return &Handler{service: service}
}

type APIShortenRequestData struct {
	URL string `json:"url" validate:"required"`
}

type APIShortenResponseData struct {
	Result string `json:"result"`
}

type APIShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	URL           string `json:"original_url" validate:"required"`
}

type APIShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"short_url"`
}

type APIShortenBatchRequestWithMask struct {
	APIShortenBatchRequest
	Mask string
}

type APIGetByUserIDResponse struct {
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
}
