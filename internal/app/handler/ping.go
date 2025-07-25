package handler

import (
	"net/http"
)

// Ping check database connection
func (h *Handler) Ping(res http.ResponseWriter, req *http.Request) {
	err := h.service.Ping()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}
