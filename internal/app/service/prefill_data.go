package service

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/repository"
)

func prefill() error {
	var err error
	if config.DatabaseDSN != "" {
		err = prefillFromDB()
		if err != nil {
			return err
		}
	} else {
		err = prefillFromFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func prefillFromFile() error {
	err := repository.DBHandler.OpenReadOnly()
	defer repository.DBHandler.Close()
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

func prefillFromDB() error {
	results, err := repository.GetAllRecordsFromDB()
	if err != nil {
		return err
	}
	for _, result := range results {
		repository.CreateInMemory(result.Id, result.OriginalURL)
	}
	return nil
}
