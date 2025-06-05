package utils

import (
	"math/rand"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func CreateRandomMask() string {
	result := make([]byte, config.ShortURLLength)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
