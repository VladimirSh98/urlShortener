package handler

import (
	"context"
	"encoding/json"
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	myProto "github.com/VladimirSh98/urlShortener/proto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

func TestHandleGRPCGetURLsByUser(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	testCases := []struct {
		name           string
		input          []getByUserIDResponseAPI
		expectedStatus int
		expectError    bool
	}{
		{
			name: "success case with multiple URLs",
			input: []getByUserIDResponseAPI{
				{ShortURL: "http://short/abc", URL: "http://original.com/long1"},
				{ShortURL: "http://short/def", URL: "http://original.com/long2"},
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "empty response",
			input:          []getByUserIDResponseAPI{},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			h := &Handler{}
			h.handleGRPCGetURLsByUser(rr, tc.input, sugar)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
			if !tc.expectError && rr.Code == http.StatusOK {
				var response myProto.GetUserURLsResponse
				if err := proto.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(response.Urls) != len(tc.input) {
					t.Errorf("expected %d URLs, got %d", len(tc.input), len(response.Urls))
				}
				for i, url := range response.Urls {
					if url.ShortUrl != tc.input[i].ShortURL || url.OriginalUrl != tc.input[i].URL {
						t.Errorf("URL mismatch at index %d", i)
					}
				}
				if ct := rr.Header().Get("Content-Type"); ct != "application/grpc+proto" {
					t.Errorf("expected Content-Type 'application/grpc+proto', got '%s'", ct)
				}
			}
		})
	}
}
