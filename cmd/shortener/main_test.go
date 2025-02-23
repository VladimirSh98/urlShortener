package main

import (
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
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			request := httptest.NewRequest(test.testRequest.method, test.testRequest.URL, strings.NewReader(test.testRequest.body))
			w := httptest.NewRecorder()
			createShortURL(w, request)
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

func TestReturnFullURL(t *testing.T) {
	
}
