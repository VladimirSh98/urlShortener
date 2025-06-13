package handler

import (
	"context"
	"encoding/json"
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	shortenMock "github.com/VladimirSh98/urlShortener/mocks/shorten"
	"github.com/stretchr/testify/assert"
)

func TestGetURLsByUser(t *testing.T) {
	type expect struct {
		status   int
		response any
	}
	type testRequest struct {
		getByUserID    []dbrepo.Shorter
		errGetByUserID error
		checkResponse  bool
	}
	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
	}{
		{
			description: "Test #1. Database error",
			expect: expect{
				status: http.StatusBadRequest,
			},
			testRequest: testRequest{
				errGetByUserID: errors.New("database error"),
				checkResponse:  false,
			},
		},
		{
			description: "Test #2. No content",
			expect: expect{
				status: http.StatusNoContent,
			},
			testRequest: testRequest{
				getByUserID:    []dbrepo.Shorter{},
				errGetByUserID: nil,
				checkResponse:  false,
			},
		},
		{
			description: "Test #3. Success",
			expect: expect{
				status: http.StatusOK,
				response: []getByUserIDResponseAPI{
					{
						ShortURL: "/ffsdafd",
						URL:      "http://test.com",
					},
				},
			},
			testRequest: testRequest{
				getByUserID: []dbrepo.Shorter{
					{
						ID:          "ffsdafd",
						OriginalURL: "http://test.com",
						UserID:      1,
						Archived:    false,
					},
				},
				errGetByUserID: nil,
				checkResponse:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodDelete, "/api/user/urls", nil,
			)
			ctx := context.WithValue(request.Context(), middleware.UserIDKey, 1)
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := shortenMock.NewMockShortenServiceInterface(ctrl)
			mockService.EXPECT().GetByUserID(gomock.Any()).Return(
				test.testRequest.getByUserID, test.testRequest.errGetByUserID).AnyTimes()
			mockHandler := Handler{service: mockService}
			mockHandler.ManagerGetURLsByUser(w, request.WithContext(ctx))
			result := w.Result()
			assert.Equal(t, test.expect.status, result.StatusCode, "Неверный код ответа")
			defer result.Body.Close()
			if test.testRequest.checkResponse {
				body, err := io.ReadAll(result.Body)
				var realResponse []getByUserIDResponseAPI
				assert.NoError(t, err)
				err = json.Unmarshal(body, &realResponse)
				assert.NoError(t, err)
				assert.Equal(t, realResponse, test.expect.response)
			}
		})
	}
}
