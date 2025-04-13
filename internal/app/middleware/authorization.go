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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: UserCount,
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", 0, err
	}
	return tokenString, UserCount, nil
}

func GetUserID(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return 0, customErr.ParseTokenError
	}
	if !token.Valid {
		return 0, customErr.NotValidTokenError
	}
	return claims.UserID, nil
}

func Authorize(request *http.Request) error {
	var err error
	var cookie *http.Cookie
	var userID int
	cookie, err = request.Cookie("Authorization")
	if errors.Is(err, http.ErrNoCookie) {
		_, userID, err = BuildJWTString()
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	} else {
		userID, err = GetUserID(cookie.Value)
		if errors.Is(err, customErr.ParseTokenError) || errors.Is(err, customErr.NotValidTokenError) {
			_, userID, err = BuildJWTString()
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	stringUserID := fmt.Sprintf("%d", userID)
	request.AddCookie(&http.Cookie{Name: "userID", Value: stringUserID})
	return nil
}
