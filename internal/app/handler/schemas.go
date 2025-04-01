package handler

type APIShortenRequestData struct {
	URL string `json:"url" validate:"required"`
}

type APIShortenResponseData struct {
	Result string `json:"result"`
}
