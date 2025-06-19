package database

// Ping check database connection
func (repo *ShortenRepository) Ping() error {
	err := repo.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}
