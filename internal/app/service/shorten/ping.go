package shorten

// Ping check database connection
func (s ShortenService) Ping() error {
	err := s.Repo.Ping()
	if err != nil {
		return err
	}
	return nil
}
