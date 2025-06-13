package file

// Close file
func (handler *handler) Close() error {
	err := handler.file.Close()
	if err != nil {
		return err
	}
	return nil
}
