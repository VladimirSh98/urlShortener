package shorten

import (
	dbrepo "github.com/VladimirSh98/urlShortener/internal/app/repository/database"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/file"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
	"go.uber.org/zap"
)

// BatchCreate batch create urls
func (s ShortenService) BatchCreate(data []dbrepo.ShortenBatchRequest) {
	sugar := zap.S()
	err := s.Repo.BatchCreate(data)
	if err != nil {
		sugar.Warnln("Failed write to database %v", err)
	}
	go func() {
		for _, record := range data {
			memory.CreateInMemory(record.Mask, record.URL)
			err = file.CreateInFile(record.Mask, record.URL)
			if err != nil {
				sugar.Warnln("Failed write to file %v", err)
			}
		}
	}()
}
