package middleware

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckTrustedSubnet(t *testing.T) {
	originalSubnet := config.TrustedSubnet
	defer func() { config.TrustedSubnet = originalSubnet }()

	tests := []struct {
		name           string
		trustedSubnet  string
		xRealIP        string
		expectedStatus int
	}{
		{
			name:           "Empty trusted subnet",
			trustedSubnet:  "",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Missing X-Real-IP header",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "IP not in trusted subnet",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "10.0.0.1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "IP in trusted subnet",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.100",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid trusted subnet",
			trustedSubnet:  "invalid_subnet",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.TrustedSubnet = tt.trustedSubnet

			req := httptest.NewRequest("GET", "http://example.com/api/internal/stats", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			rr := httptest.NewRecorder()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := CheckTrustedSubnet(mockHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
