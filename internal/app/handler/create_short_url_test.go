package handler

import (
	"context"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	shortenMock "github.com/VladimirSh98/urlShortener/mocks/shorten"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateShortURL(t *testing.T) {
	type expect struct {
		status          int
		contentType     string
		checkBodyLength bool
	}
	type testRequest struct {
		URL    string
		method string
		body   string
		err    error
	}
	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
	}{
		{
			description: "Test #1. Wrong request",
			expect: expect{
				status:          http.StatusBadRequest,
				contentType:     "",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				URL:    "/",
				method: http.MethodPatch,
				body:   "",
			},
		},
		{
			description: "Test #2. Wrong request",
			expect: expect{
				status:          http.StatusBadRequest,
				contentType:     "",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				URL:    "/qwe",
				method: http.MethodPost,
				body:   "",
			},
		},
		{
			description: "Test #3. Wrong body",
			expect: expect{
				status:          http.StatusBadRequest,
				contentType:     "",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				URL:    "/",
				method: http.MethodPost,
				body:   "",
			},
		},
		{
			description: "Test #4. Success",
			expect: expect{
				status:          http.StatusCreated,
				contentType:     "text/plain",
				checkBodyLength: true,
			},
			testRequest: testRequest{
				URL:    "/",
				method: http.MethodPost,
				body:   "http://example.com",
			},
		},
		{
			description: "Test #5. Error",
			expect: expect{
				status:          http.StatusConflict,
				contentType:     "text/plain",
				checkBodyLength: true,
			},
			testRequest: testRequest{
				URL:    "/",
				method: http.MethodPost,
				body:   "http://example.com",
				err:    customErr.ErrConstraintViolation,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(test.testRequest.method, test.testRequest.URL, strings.NewReader(test.testRequest.body))
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := shortenMock.NewMockShortenServiceInterface(ctrl)
			mockService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(
				"", test.testRequest.err).AnyTimes()
			ctx := context.WithValue(request.Context(), middleware.UserIDKey, 1)
			customHandler := NewHandler(mockService)
			customHandler.ManagerCreateShortURL(w, request.WithContext(ctx))
			result := w.Result()
			assert.Equal(t, test.expect.status, result.StatusCode, "Неверный код ответа")
			defer result.Body.Close()
			body, err := io.ReadAll(result.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.expect.contentType, result.Header.Get("Content-Type"), "Неверный тип контента в хедере")
			if len(body) == 0 && test.expect.checkBodyLength {
				t.Errorf("Отсутствует тело ответа")
			}
		})
	}
}
