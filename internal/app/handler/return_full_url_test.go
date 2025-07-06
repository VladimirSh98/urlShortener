package handler

import (
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupGlobalURLStorageCase() func() {
	memory.CreateInMemory("TestCase", "http://example.com")
	return func() {
		memory.Delete("TestCase")
	}
}

func TestReturnFullURL(t *testing.T) {
	type expect struct {
		status   int
		location string
	}
	tests := []struct {
		description string
		URL         string
		expect      expect
	}{
		{
			description: "Test #1. Not exist link",
			URL:         "/eqwrewerw",
			expect: expect{
				status:   http.StatusBadRequest,
				location: "",
			},
		},
		{
			description: "Test #2. Success case",
			URL:         "/TestCase",
			expect: expect{
				status:   http.StatusTemporaryRedirect,
				location: "http://example.com",
			},
		},
	}
	setup := setupGlobalURLStorageCase()
	defer setup()
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.URL, nil)
			request.SetPathValue("id", test.URL[1:])
			w := httptest.NewRecorder()
			repo := dbRepo.NewShortenRepository(database.DBConnection.Conn)
			service := shorten.NewShortenService(repo)
			customHandler := NewHandler(service)
			customHandler.ManagerReturnFullURL(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, test.expect.status, result.StatusCode)
			if test.expect.location != "" {
				assert.Equal(t, test.expect.location, result.Header.Get("Location"))
			}
		})
	}
}
