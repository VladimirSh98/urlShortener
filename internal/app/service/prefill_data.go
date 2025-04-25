package service

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"github.com/VladimirSh98/urlShortener/internal/app/database"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	dbRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	fileRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/file"
	memoryRepo "github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
	"github.com/VladimirSh98/urlShortener/internal/app/service/shorten"
)

func Prefill() error {
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
	err := fileRepo.DBHandler.OpenReadOnly()
	defer fileRepo.DBHandler.Close()
	if err != nil {
		return err
	}
	var record *fileRepo.URLStorageFileData
	for {
		record, err = fileRepo.DBHandler.ReadLine()
		if err != nil {
			return err
		}
		if record == nil {
			return nil
		}
		fileRepo.DBHandler.Count++
		memoryRepo.CreateInMemory(record.ShortURL, record.OriginalURL)
	}
}

func prefillFromDB() error {
	getService := shorten.NewShortenService(dbRepo.ShortenRepository{Conn: database.DBConnection.Conn})
	results, err := getService.GetAllRecords()
	if err != nil {
		return err
	}
	for _, result := range results {
		memoryRepo.CreateInMemory(result.ID, result.OriginalURL)
		if int(middleware.UserCount) < result.UserID {
			middleware.UserCount = int64(result.UserID)
		}
	}
	return nil
}
