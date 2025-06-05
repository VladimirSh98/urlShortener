package middleware

import (
	"net/http"
	"sync/atomic"
	"time"

	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

func BuildJWTString() (string, int, error) {
	atomic.AddInt64(&UserCount, 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: UserCount,
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", 0, err
	}
	return tokenString, int(UserCount), nil
}

func GetUserID(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0, customErr.ErrParseToken
	}
	if !token.Valid {
		return 0, customErr.ErrNotValidToken
	}
	return int(claims.UserID), nil
}

func Authorize(request *http.Request) (string, int, error) {
	var err error
	var cookie *http.Cookie
	var userID int
	var token string

	cookie, err = request.Cookie("Authorization")
	if errors.Is(err, http.ErrNoCookie) {
		token, userID, err = BuildJWTString()
		if err != nil {
			return "", 0, err
		}

	} else if err != nil {
		return "", 0, err
	} else {
		token = cookie.Value
		userID, err = GetUserID(token)
		if errors.Is(err, customErr.ErrParseToken) || errors.Is(err, customErr.ErrNotValidToken) {
			token, userID, err = BuildJWTString()
			if err != nil {
				return "", 0, err
			}
		} else if err != nil {
			return "", 0, err
		}
	}

	return token, userID, nil
}
