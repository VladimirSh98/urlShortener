package database

import "fmt"

// BatchCreate batch create urls
func (repo *ShortenRepository) BatchCreate(data []ShortenBatchRequest) error {
	queries := make([]string, 0)
	for _, record := range data {
		query := fmt.Sprintf(
			"INSERT INTO urls (id, original_url, user_id, archived) VALUES ('%s', '%s', '%d', false);", record.Mask, record.URL, record.UserID)
		queries = append(queries, query)
	}
	err := batchCreate(repo, queries)
	if err != nil {
		return err
	}
	return nil
}

func batchCreate(repo *ShortenRepository, queries []string) error {
	tx, err := repo.Conn.Begin()
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
