package middleware

import (
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"net/http"
	"sync/atomic"
)

func buildJWTString() (string, int, error) {
	atomic.AddInt64(&UserCount, 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           UserCount,
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", 0, err
	}
	return tokenString, int(UserCount), nil
}

func getUserID(tokenString string) (int, error) {
	userClaims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, userClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return 0, customErr.ErrParseToken
	}
	if !token.Valid {
		return 0, customErr.ErrNotValidToken
	}
	return int(userClaims.UserID), nil
}

func authorize(request *http.Request) (string, int, error) {
	var err error
	var cookie *http.Cookie
	var userID int
	var token string

	cookie, err = request.Cookie("Authorization")
	if errors.Is(err, http.ErrNoCookie) {
		token, userID, err = buildJWTString()
		if err != nil {
			return "", 0, err
		}

	} else if err != nil {
		return "", 0, err
	} else {
		token = cookie.Value
		userID, err = getUserID(token)
		if errors.Is(err, customErr.ErrParseToken) || errors.Is(err, customErr.ErrNotValidToken) {
			token, userID, err = buildJWTString()
			if err != nil {
				return "", 0, err
			}
		} else if err != nil {
			return "", 0, err
		}
	}

	return token, userID, nil
}
