package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
)

type ShortenServiceInterface interface {
	Create(mask string, originalURL string, userID int) (string, error)
	GetAllRecords() ([]dbrepo.Shorter, error)
	BatchCreate(data []dbrepo.ShortenBatchRequest)
	GetByOriginalURL(originalURL string) (dbrepo.Shorter, error)
	GetByUserID(userID int) ([]dbrepo.Shorter, error)
	GetByShortURL(shortURL string) (dbrepo.Shorter, error)
	BatchUpdate(data []string, userID int) error
	Ping() error
}

type ShortenService struct {
	Repo dbrepo.ShortenRepository
}

func NewShortenService(repo dbrepo.ShortenRepository) ShortenServiceInterface {
	return &ShortenService{Repo: repo}
}
