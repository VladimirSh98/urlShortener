package file

func (handler *Handler) Close() error {
	err := handler.file.Close()
	if err != nil {
		return err
	}
	return nil
}
