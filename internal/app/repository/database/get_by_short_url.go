package database

import "fmt"

func (repo *ShortenRepository) GetByShortURL(shortURL string) (Shorter, error) {
	if repo.Conn == nil {
		return Shorter{}, fmt.Errorf("DB connection is not open")
	}
	var record Shorter
	query := fmt.Sprintf("SELECT * FROM urls WHERE id = '%s' limit 1", shortURL)
	row := repo.Conn.QueryRow(query)
	err := row.Scan(&record.ID, &record.OriginalURL, &record.CreatedAt, &record.UserID, &record.Archived)
	if err != nil {
		return record, err
	}
	return record, nil
}
