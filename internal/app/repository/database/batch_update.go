package database

import "github.com/lib/pq"

// BatchUpdate archived urls by ID and user ID
func (repo *ShortenRepository) BatchUpdate(data []string, userID int) error {
	query := "UPDATE urls SET archived = true WHERE id = ANY($1::text[]) and user_id = $2"
	_, err := repo.Conn.Exec(query, pq.Array(data), userID)
	if err != nil {
		return err
	}
	return nil
}
