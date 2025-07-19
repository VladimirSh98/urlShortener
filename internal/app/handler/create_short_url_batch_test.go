package handler

import (
	"bytes"
	"context"
	"encoding/json"
	shortenMock "github.com/VladimirSh98/urlShortener/mocks/shorten"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortURLBatch(t *testing.T) {
	type expect struct {
		status          int
		contentType     string
		checkBodyLength bool
	}
	type testRequest struct {
		body    any
		headers map[string]string
	}
	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
	}{
		{
			description: "Test #1. Zero request body",
			expect: expect{
				status:          http.StatusCreated,
				contentType:     "application/json",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				body: []struct{}{},
			},
		},
		{
			description: "Test #2. Not valid request body",
			expect: expect{
				status:          http.StatusBadRequest,
				contentType:     "",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				body: struct {
					Curl string `json:"curl"`
				}{
					Curl: "http://example.com",
				},
			},
		},
		{
			description: "Test #3. Not valid request body",
			expect: expect{
				status:          http.StatusBadRequest,
				contentType:     "",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				body: "ffefefewfwe",
			},
		},
		{
			description: "Test #4. Success",
			expect: expect{
				status:          http.StatusCreated,
				contentType:     "application/json",
				checkBodyLength: true,
			},
			testRequest: testRequest{
				body: []shortenBatchRequestAPI{
					{URL: "http://example.com", CorrelationID: "123"},
				},
			},
		},
		{
			description: "Test #5. GRPC Success",
			expect: expect{
				status:          http.StatusOK,
				contentType:     "application/grpc+proto",
				checkBodyLength: true,
			},
			testRequest: testRequest{
				body: []shortenBatchRequestAPI{
					{URL: "http://example.com", CorrelationID: "123"},
				},
				headers: map[string]string{
					"Accept": "application/grpc",
				},
			},
		},
		{
			description: "Test #6. GRPC Invalid Body",
			expect: expect{
				status:          http.StatusBadRequest,
				contentType:     "",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				body: "invalid body",
				headers: map[string]string{
					"Accept": "application/grpc",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			jsonBody, _ := json.Marshal(test.testRequest.body)
			request := httptest.NewRequest(
				http.MethodPost, "/api/shorten", bytes.NewReader(jsonBody),
			)
			if test.testRequest.headers != nil {
				for k, v := range test.testRequest.headers {
					request.Header.Set(k, v)
				}
			}
			ctx := context.WithValue(request.Context(), middleware.UserIDKey, 1)
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := shortenMock.NewMockShortenServiceInterface(ctrl)
			service.EXPECT().BatchCreate(gomock.Any()).Return().AnyTimes()
			customHandler := NewHandler(service)
			customHandler.ManagerCreateShortURLBatch(w, request.WithContext(ctx))
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
