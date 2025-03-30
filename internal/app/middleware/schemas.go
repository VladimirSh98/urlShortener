package middleware

import (
	"io"
	"net/http"
	"sync"
)

type customResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
	once   sync.Once
}

type compressWriter struct {
	customResponseWriter
	Writer io.Writer
	once   sync.Once
}
