package database

func (repo *ShortenRepository) GetAllRecords() ([]Shorter, error) {
	query := "SELECT * FROM urls"
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
