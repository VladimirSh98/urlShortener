package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
)

// GetByOriginalURL get short URL by original URL
func (s ShortenService) GetByOriginalURL(originalURL string) (dbrepo.Shorter, error) {
	records, err := s.Repo.GetByOriginalURL(originalURL)
	if err != nil {
		return dbrepo.Shorter{}, err
	}
	return records, nil
}
