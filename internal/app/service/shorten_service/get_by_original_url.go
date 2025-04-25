package shorten_service

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
)

func (s ShortenService) GetByOriginalURL(originalURL string) (dbrepo.Shorter, error) {
	records, err := s.Repo.GetByOriginalURL(originalURL)
	if err != nil {
		return dbrepo.Shorter{}, err
	}
	return records, nil
}
