package handler

import (
	"encoding/json"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	shortenMock "github.com/VladimirSh98/urlShortener/mocks/shorten"
	myProto "github.com/VladimirSh98/urlShortener/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStats(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := shortenMock.NewMockShortenServiceInterface(ctrl)

	h := &Handler{service: mockService}

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   getStatsResponseAPI
		headers        map[string]string
	}{
		{
			name: "Success case",
			mockSetup: func() {
				mockService.EXPECT().GetAllRecords().Return([]dbrepo.Shorter{
					{ID: "eewe", OriginalURL: "http://example.com/1", UserID: 1},
					{ID: "eewe1", OriginalURL: "http://example.com/2", UserID: 1},
					{ID: "eewe2", OriginalURL: "http://example.com/3", UserID: 2},
				}, nil)
				middleware.UserCount = 2
			},
			expectedStatus: http.StatusOK,
			expectedBody: getStatsResponseAPI{
				URLS:  3,
				Users: 2,
			},
		},
		{
			name: "Service error case",
			mockSetup: func() {
				mockService.EXPECT().GetAllRecords().Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "GRPC Success case",
			mockSetup: func() {
				mockService.EXPECT().GetAllRecords().Return([]dbrepo.Shorter{
					{ID: "eewe", OriginalURL: "http://example.com/1", UserID: 1},
					{ID: "eewe1", OriginalURL: "http://example.com/2", UserID: 1},
					{ID: "eewe2", OriginalURL: "http://example.com/3", UserID: 2},
				}, nil)
				middleware.UserCount = 2
			},
			headers: map[string]string{
				"Accept": "application/grpc",
			},
			expectedStatus: http.StatusOK,
			expectedBody: getStatsResponseAPI{
				URLS:  3,
				Users: 2,
			},
		},
		{
			name: "GRPC Service error case",
			mockSetup: func() {
				mockService.EXPECT().GetAllRecords().Return(nil, assert.AnError)
			},
			headers: map[string]string{
				"Accept": "application/grpc",
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
			if tt.headers != nil {
				for k, v := range tt.headers {
					req.Header.Set(k, v)
				}
			}
			rr := httptest.NewRecorder()
			h.GetStats(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.headers != nil {
				if tt.expectedStatus == http.StatusOK {
					var response myProto.StatsResponse
					err := proto.Unmarshal(rr.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Equal(t, int32(tt.expectedBody.URLS), response.Urls)
					assert.Equal(t, int32(tt.expectedBody.Users), response.Users)
				}
			} else {
				if tt.expectedStatus == http.StatusOK {
					var response getStatsResponseAPI
					err := json.NewDecoder(rr.Body).Decode(&response)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedBody, response)
				}
			}
		})
	}
}
