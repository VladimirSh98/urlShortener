package shorten

// BatchUpdate archived urls by ID and user ID
func (s ShortenService) BatchUpdate(data []string, userID int) error {
	err := s.Repo.BatchUpdate(data, userID)
	if err != nil {
		return err
	}
	return nil
}
