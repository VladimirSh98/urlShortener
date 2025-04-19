package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"sync"
	"time"
)

var UserCount int

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"

type customResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
	once   sync.Once
}

type compressWriter struct {
	http.ResponseWriter
	size   int
	status int
	once   sync.Once
	Writer io.Writer
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}
