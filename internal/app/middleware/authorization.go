package middleware

import (
	"fmt"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func BuildJWTString() (string, int, error) {
	UserCount++
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
	return tokenString, UserCount, nil
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
	return claims.UserID, nil
}

func Authorize(request *http.Request) (string, error) {
	var err error
	var cookie *http.Cookie
	var userID int
	var token string

	cookie, err = request.Cookie("Authorization")
	if errors.Is(err, http.ErrNoCookie) {
		token, userID, err = BuildJWTString()
		if err != nil {
			return "", err
		}

	} else if err != nil {
		return "", err
	} else {
		token = cookie.Value
		userID, err = GetUserID(token)
		if errors.Is(err, customErr.ErrParseToken) || errors.Is(err, customErr.ErrNotValidToken) {
			token, userID, err = BuildJWTString()
			if err != nil {
				return "", err
			}
		} else if err != nil {
			return "", err
		}
	}

	stringUserID := fmt.Sprintf("%d", userID)
	request.AddCookie(&http.Cookie{Name: "userID", Value: stringUserID})
	return token, nil
}
