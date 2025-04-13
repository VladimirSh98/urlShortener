package repository

import (
	"database/sql"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
)

func createDB(mask string, originalURL string, userID int) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO urls (id, original_url, user_id) VALUES ('%s', '%s', '%d');", mask, originalURL, userID)
	res, err := database.DBConnection.Exec(query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetAllRecordsFromDB() ([]Shortner, error) {
	query := "SELECT * FROM urls"
	rows, err := database.DBConnection.Query(query)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	results := make([]Shortner, 0)
	for rows.Next() {
		var record Shortner
		err = rows.Scan(&record.ID, &record.OriginalURL, &record.CreatedAt, &record.UserID)
		if err != nil {
			return nil, err
		}

		results = append(results, record)
	}
	return results, nil
}

func BatchCreateDB(data []ShortenBatchRequest) error {
	queries := make([]string, 0)
	for _, record := range data {
		query := fmt.Sprintf("INSERT INTO urls (id, original_url, user_id) VALUES ('%s', '%s', '%d');", record.Mask, record.URL, record.UserID)
		queries = append(queries, query)
	}
	err := database.DBConnection.BatchCreate(queries)
	if err != nil {
		return err
	}
	return nil
}

func GetByOriginalURLFromBD(originalURL string) (Shortner, error) {
	query := fmt.Sprintf("SELECT * FROM urls WHERE original_url = '%s' limit 1", originalURL)
	row := database.DBConnection.QueryRow(query)
	var record Shortner
	err := row.Scan(&record.ID, &record.OriginalURL, &record.CreatedAt, &record.UserID)
	if err != nil {
		return record, err
	}
	return record, nil
}

func GetByUserID(userID int) ([]Shortner, error) {
	query := fmt.Sprintf("SELECT * FROM urls WHERE user_id = '%d'", userID)
	rows, err := database.DBConnection.Query(query)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	results := make([]Shortner, 0)
	for rows.Next() {
		var record Shortner
		err = rows.Scan(&record.ID, &record.OriginalURL, &record.CreatedAt, &record.UserID)
		if err != nil {
			return nil, err
		}

		results = append(results, record)
	}
	return results, nil
}
