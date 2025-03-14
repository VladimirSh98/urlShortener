package middleware

import "net/http"

type customResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}
