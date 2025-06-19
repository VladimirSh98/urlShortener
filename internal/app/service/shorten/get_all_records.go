package shorten

import dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"

// GetAllRecords get all records from urls
func (s ShortenService) GetAllRecords() ([]dbrepo.Shorter, error) {
	records, err := s.Repo.GetAllRecords()
	if err != nil {
		return nil, err
	}
	return records, nil
}
