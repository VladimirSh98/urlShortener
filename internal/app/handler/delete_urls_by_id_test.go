package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	shortenMock "github.com/VladimirSh98/urlShortener/mocks/shorten"
	"github.com/stretchr/testify/assert"
)

func TestDeleteURLsByID(t *testing.T) {
	type expect struct {
		status int
	}
	type testRequest struct {
		body any
	}
	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
	}{
		{
			description: "Test #1. Zero request body",
			expect: expect{
				status: http.StatusBadRequest,
			},
			testRequest: testRequest{
				body: "",
			},
		},
		{
			description: "Test #2. Wrong request body",
			expect: expect{
				status: http.StatusBadRequest,
			},
			testRequest: testRequest{
				body: "fdfdsfdsf",
			},
		},
		{
			description: "Test #3. Success",
			expect: expect{
				status: http.StatusAccepted,
			},
			testRequest: testRequest{
				body: []string{"fdfdsfdsf"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			requestBody, err := json.Marshal(test.testRequest.body)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}
			request := httptest.NewRequest(
				http.MethodDelete, "/api/user/urls", bytes.NewBuffer(requestBody),
			)
			ctx := context.WithValue(request.Context(), middleware.UserIDKey, 1)
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := shortenMock.NewMockShortenServiceInterface(ctrl)
			mockService.EXPECT().BatchUpdate(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			mockHandler := Handler{service: mockService}
			mockHandler.ManagerDeleteURLsByID(w, request.WithContext(ctx))
			result := w.Result()
			assert.Equal(t, test.expect.status, result.StatusCode, "Неверный код ответа")
			defer result.Body.Close()
		})
	}
}
