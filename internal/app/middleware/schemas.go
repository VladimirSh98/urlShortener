package middleware

import (
	"io"
	"net/http"
)

type customResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

type compressWriter struct {
	customResponseWriter
	Writer io.Writer
}
