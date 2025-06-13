package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"sync"
)

// UserCount count user with shorten urls
var UserCount int64

type contextKey string

// UserIDKey contains contextKey for userID
const UserIDKey contextKey = "userID"

// SecretKey contains secret key for tokens
const SecretKey = "supersecretkey"

type customResponseWriter struct {
	http.ResponseWriter
	Size   int
	Status int
	once   sync.Once
}

type compressWriter struct {
	http.ResponseWriter
	Size   int
	Status int
	once   sync.Once
	Writer io.Writer
}

type claims struct {
	jwt.RegisteredClaims
	UserID int64
}

// WriteHeader add status code to response
func (lrw *customResponseWriter) WriteHeader(code int) {
	lrw.Status = code
	lrw.once.Do(func() { lrw.ResponseWriter.WriteHeader(code) })
}

// Write add size to response
func (lrw *customResponseWriter) Write(body []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(body)
	lrw.Size += n
	return n, err
}

// WriteHeader add status code to response
func (lrw *compressWriter) WriteHeader(code int) {
	lrw.Status = code
	lrw.once.Do(func() { lrw.ResponseWriter.WriteHeader(code) })
}

// Write add size to response
func (lrw *compressWriter) Write(body []byte) (int, error) {
	n, err := lrw.Writer.Write(body)
	lrw.Size += n
	return n, err
}
