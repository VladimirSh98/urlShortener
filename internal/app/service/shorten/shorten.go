package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
)

// ShortenServiceInterface interface for service
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

// ShortenService service with ShortenRepository
type ShortenService struct {
	Repo dbrepo.ShortenRepoInterface
}

// NewShortenService create new service with ShortenRepository
func NewShortenService(repo dbrepo.ShortenRepoInterface) ShortenServiceInterface {
	return &ShortenService{Repo: repo}
}
