package repository

import (
	"database/sql"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
)

func createDB(mask string, originalURL string) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO urls (id, original_url) VALUES ('%s', '%s');", mask, originalURL)
	res, err := database.DBConnection.Exec(query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetAllRecordsFromDB() ([]Shortner, error) {
	query := "SELECT * FROM urls"
	results := make([]Shortner, 0)
	rows, err := database.DBConnection.Query(query)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var record Shortner
		err = rows.Scan(&record.ID, &record.OriginalURL, &record.CreatedAt)
		if err != nil {
			return nil, err
		}

		results = append(results, record)
	}
	return results, nil
}
