package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"sync"
	"time"
)

var UserCount int64

type contextKey string

const UserIDKey contextKey = "userID"

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

type CustomResponseWriter struct {
	http.ResponseWriter
	Size   int
	Status int
	once   sync.Once
}

type CompressWriter struct {
	http.ResponseWriter
	Size   int
	Status int
	once   sync.Once
	Writer io.Writer
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

func (lrw *CustomResponseWriter) WriteHeader(code int) {
	lrw.Status = code
	lrw.once.Do(func() { lrw.ResponseWriter.WriteHeader(code) })
}

func (lrw *CustomResponseWriter) Write(body []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(body)
	lrw.Size += n
	return n, err
}

func (lrw *CompressWriter) WriteHeader(code int) {
	lrw.Status = code
	lrw.once.Do(func() { lrw.ResponseWriter.WriteHeader(code) })
}

func (lrw *CompressWriter) Write(body []byte) (int, error) {
	n, err := lrw.Writer.Write(body)
	lrw.Size += n
	return n, err
}
