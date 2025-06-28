package handler

import (
	shortenMock "github.com/VladimirSh98/urlShortener/mocks/shorten"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	type expect struct {
		status int
	}
	type testRequest struct {
		err any
	}
	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
	}{
		{
			description: "Test #1. Error Ping",
			expect: expect{
				status: http.StatusInternalServerError,
			},
			testRequest: testRequest{
				err: errors.New("request body cannot be empty"),
			},
		},
		{
			description: "Test #2. Success Ping",
			expect: expect{
				status: http.StatusOK,
			},
			testRequest: testRequest{
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodGet, "/ping", nil,
			)
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := shortenMock.NewMockShortenServiceInterface(ctrl)
			mockService.EXPECT().Ping().Return(test.testRequest.err).AnyTimes()
			mockHandler := Handler{service: mockService}
			mockHandler.Ping(w, request)
			result := w.Result()
			assert.Equal(t, test.expect.status, result.StatusCode, "Неверный код ответа")
		})
	}
}
