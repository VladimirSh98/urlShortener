package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
)

// GetByShortURL get original URL by short URL
func (s ShortenService) GetByShortURL(shortURL string) (dbrepo.Shorter, error) {
	records, err := s.Repo.GetByShortURL(shortURL)
	if err != nil {
		return dbrepo.Shorter{}, err
	}
	return records, nil
}
