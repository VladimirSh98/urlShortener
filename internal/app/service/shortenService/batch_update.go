package shortenService

func (s ShortenService) BatchUpdate(data []string, userID int) error {
	err := s.Repo.BatchUpdate(data, userID)
	if err != nil {
		return err
	}
	return nil
}
