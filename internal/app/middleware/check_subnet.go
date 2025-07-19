package middleware

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"net"
	"net/http"
)

// CheckTrustedSubnet check trusted subnet
func CheckTrustedSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.TrustedSubnet == "" {
			http.Error(w, "Access forbidden", http.StatusForbidden)
			return
		}

		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			http.Error(w, "X-Real-IP header required", http.StatusForbidden)
			return
		}

		_, subnet, err := net.ParseCIDR(config.TrustedSubnet)
		if err != nil {
			http.Error(w, "Invalid trusted subnet configuration", http.StatusInternalServerError)
			return
		}

		ip := net.ParseIP(realIP)
		if ip == nil || !subnet.Contains(ip) {
			http.Error(w, "Access forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
