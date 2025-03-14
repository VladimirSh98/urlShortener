package handler

type ApiShortenRequestData struct {
	URL string `json:"url" validate:"required"`
}

type ApiShortenResponse struct {
	Result string `json:"result"`
}
