package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
)

func BenchmarkManagerCreateShortURL(b *testing.B) {
	for i := 0; i < 20; i++ {
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		w := httptest.NewRecorder()
		ctx := context.WithValue(request.Context(), middleware.UserIDKey, 1)
		repo := dbRepo.ShortenRepository{Conn: database.DBConnection.Conn}
		service := shorten.NewShortenService(repo)
		customHandler := NewHandler(service)
		customHandler.ManagerCreateShortURL(w, request.WithContext(ctx))
	}
}
