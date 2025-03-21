package service

import "github.com/VladimirSh98/urlShortener/internal/app/repository"

func prefill() error {
	err := repository.DBHandler.OpenReadOnly()
	if err != nil {
		return err
	}
	var record *repository.URLStorageFileData
	for {
		record, err = repository.DBHandler.ReadLine()
		if err != nil {
			return err
		}
		if record == nil {
			return nil
		}
		repository.DBHandler.Count++
		repository.Create(record.ShortURL, record.OriginalURL, false)
	}
}
