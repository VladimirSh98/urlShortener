package database

import "fmt"

// GetByUserID get user URL by user ID from database
func (repo *ShortenRepository) GetByUserID(userID int) ([]Shorter, error) {
	query := fmt.Sprintf("SELECT * FROM urls WHERE user_id = '%d'", userID)
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	results := make([]Shorter, 0)
	for rows.Next() {
		var record Shorter
		err = rows.Scan(&record.ID, &record.OriginalURL, &record.CreatedAt, &record.UserID, &record.Archived)
		if err != nil {
			return nil, err
		}

		results = append(results, record)
	}
	return results, nil
}
