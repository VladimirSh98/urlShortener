package shorten_service

func (s ShortenService) Ping() error {
	err := s.Repo.Ping()
	if err != nil {
		return err
	}
	return nil
}
