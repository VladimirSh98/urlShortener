package database

import (
	"database/sql"
	"fmt"
)

// Create record in urls by mask, original URL and user ID
func (repo *ShortenRepository) Create(mask string, originalURL string, userID int) (sql.Result, error) {
	if repo.Conn == nil {
		return nil, fmt.Errorf("DB connection is not open")
	}
	query := fmt.Sprintf("INSERT INTO urls (id, original_url, user_id, archived) VALUES ('%s', '%s', '%d', false);", mask, originalURL, userID)
	res, err := repo.Conn.Exec(query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
