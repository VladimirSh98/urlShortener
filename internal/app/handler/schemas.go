package handler

import "github.com/VladimirSh98/urlShortener/internal/app/service/shorten"

// Handler with ShortenServiceInterface
type Handler struct {
	service shorten.ShortenServiceInterface
}

// NewHandler create new handler
func NewHandler(service shorten.ShortenServiceInterface) *Handler {
	return &Handler{service: service}
}

type shortenRequestDataAPI struct {
	URL string `json:"url" validate:"required"`
}

type shortenResponseDataAPI struct {
	Result string `json:"result"`
}

type shortenBatchRequestAPI struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	URL           string `json:"original_url" validate:"required"`
}

type shortenBatchResponseAPI struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"short_url"`
}

type shortenBatchRequestWithMaskAPI struct {
	shortenBatchRequestAPI
	Mask string
}

type getByUserIDResponseAPI struct {
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
}
