package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
)

// GetByUserID get user URL by user ID from database
func (s ShortenService) GetByUserID(userID int) ([]dbrepo.Shorter, error) {
	records, err := s.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return records, nil
}
