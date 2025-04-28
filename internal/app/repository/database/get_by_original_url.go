package database

import "fmt"

func (repo *ShortenRepository) GetByOriginalURL(originalURL string) (Shorter, error) {
	var record Shorter
	query := fmt.Sprintf("SELECT * FROM urls WHERE original_url = '%s' limit 1", originalURL)
	row := repo.Conn.QueryRow(query)
	err := row.Scan(&record.ID, &record.OriginalURL, &record.CreatedAt, &record.UserID, &record.Archived)
	if err != nil {
		return record, err
	}
	return record, nil
}
