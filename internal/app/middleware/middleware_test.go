package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig(t *testing.T) {
	type expect struct {
		statusCode       int
		contentEncoding  string
		responseBody     string
		contextUserIDSet bool
		contextUserID    int
	}

	type testRequest struct {
		method          string
		url             string
		contentEncoding string
		acceptEncoding  string
		body            string
	}

	type authParams struct {
		token  string
		userID int
		err    error
	}

	tests := []struct {
		description string
		expect      expect
		testRequest testRequest
		auth        authParams
		handlerFunc http.HandlerFunc
	}{
		{
			description: "Test #1. Gzip request decompression",
			expect: expect{
				statusCode:       http.StatusOK,
				responseBody:     "processed: test body",
				contextUserIDSet: true,
				contextUserID:    123,
			},
			testRequest: testRequest{
				method:          "POST",
				url:             "http://example.com",
				contentEncoding: "gzip",
				body:            "test body",
			},
			auth: authParams{
				token:  "test-token",
				userID: 123,
				err:    nil,
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(r.Body)
				if err != nil {
					t.Errorf("Error reading body: %v", err)
					return
				}
				w.Write([]byte("processed: " + buf.String()))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var reqBody *bytes.Reader
			if test.testRequest.contentEncoding == "gzip" {
				var buf bytes.Buffer
				gz := gzip.NewWriter(&buf)
				_, err := gz.Write([]byte(test.testRequest.body))
				if err != nil {
					t.Fatalf("Failed to gzip request body: %v", err)
				}
				if err := gz.Close(); err != nil {
					t.Fatalf("Failed to close gzip writer: %v", err)
				}
				reqBody = bytes.NewReader(buf.Bytes())
			} else {
				reqBody = bytes.NewReader([]byte(test.testRequest.body))
			}

			req := httptest.NewRequest(test.testRequest.method, test.testRequest.url, reqBody)
			if test.testRequest.contentEncoding != "" {
				req.Header.Set("Content-Encoding", test.testRequest.contentEncoding)
			}
			if test.testRequest.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", test.testRequest.acceptEncoding)
			}

			rr := httptest.NewRecorder()

			handler := Config(test.handlerFunc)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expect.statusCode {
				t.Errorf("expected status code %d, got %d", test.expect.statusCode, rr.Code)
			}

			if ce := rr.Header().Get("Content-Encoding"); ce != test.expect.contentEncoding {
				t.Errorf("expected content encoding '%s', got '%s'", test.expect.contentEncoding, ce)
			}
			var responseBody string
			if test.expect.contentEncoding == "gzip" {
				gr, err := gzip.NewReader(rr.Body)
				if err != nil {
					t.Fatalf("Failed to create gzip reader: %v", err)
				}
				defer gr.Close()
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(gr)
				if err != nil {
					t.Fatalf("Failed to read gzipped content: %v", err)
				}
				responseBody = buf.String()
			} else {
				responseBody = rr.Body.String()
			}

			if responseBody != test.expect.responseBody {
				t.Errorf("expected response body '%s', got '%s'", test.expect.responseBody, responseBody)
			}
		})
	}
}
