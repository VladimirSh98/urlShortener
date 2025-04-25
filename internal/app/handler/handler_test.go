package handler

import (
	"bytes"
	"encoding/json"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
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
		//{
		//	description: "Test #1. Wrong request",
		//	expect: expect{
		//		status:          http.StatusBadRequest,
		//		contentType:     "",
		//		checkBodyLength: false,
		//	},
		//	testRequest: testRequest{
		//		URL:    "/",
		//		method: http.MethodPatch,
		//		body:   "",
		//	},
		//},
		//{
		//	description: "Test #2. Wrong request",
		//	expect: expect{
		//		status:          http.StatusBadRequest,
		//		contentType:     "",
		//		checkBodyLength: false,
		//	},
		//	testRequest: testRequest{
		//		URL:    "/qwe",
		//		method: http.MethodPost,
		//		body:   "",
		//	},
		//},
		//{
		//	description: "Test #3. Wrong body",
		//	expect: expect{
		//		status:          http.StatusBadRequest,
		//		contentType:     "",
		//		checkBodyLength: false,
		//	},
		//	testRequest: testRequest{
		//		URL:    "/",
		//		method: http.MethodPost,
		//		body:   "",
		//	},
		//},
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
			request.AddCookie(&http.Cookie{Name: "userID", Value: "1"})
			ManagerCreateShortURL(w, request)
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
			ManagerReturnFullURL(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, test.expect.status, result.StatusCode)
			if test.expect.location != "" {
				assert.Equal(t, test.expect.location, result.Header.Get("Location"))
			}
		})
	}
}

func TestCreateShortURLByJSON(t *testing.T) {
	type expect struct {
		status          int
		contentType     string
		checkBodyLength bool
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
				status:          http.StatusBadRequest,
				contentType:     "",
				checkBodyLength: false,
			},
			testRequest: testRequest{
				body: struct{}{},
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
			description: "Test #3. Success",
			expect: expect{
				status:          http.StatusCreated,
				contentType:     "application/json",
				checkBodyLength: true,
			},
			testRequest: testRequest{
				body: APIShortenRequestData{URL: "http://example.com"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			jsonBody, _ := json.Marshal(test.testRequest.body)
			request := httptest.NewRequest(
				http.MethodPost, "/api/shorten", bytes.NewReader(jsonBody),
			)
			request.AddCookie(&http.Cookie{Name: "userID", Value: "1"})
			w := httptest.NewRecorder()
			ManagerCreateShortURLByJSON(w, request)
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
