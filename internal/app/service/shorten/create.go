package shorten

import (
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/file"
	"github.com/VladimirSh98/urlShortener/internal/app/repository/memory"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s ShortenService) Create(mask string, originalURL string, userID int) (string, error) {
	var err error
	sugar := zap.S()
	_, err = s.Repo.Create(mask, originalURL, userID)
	if err != nil {
		sugar.Warnln("Failed write to database %v", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			sugar.Infoln("URL already exists %s", originalURL)
			var oldMask string
			oldMask, err = GetByOriginalURL(s, originalURL)
			if err != nil {
				return "", err
			}
			return oldMask, customErr.ErrConstraintViolation
		}
	}
	memory.CreateInMemory(mask, originalURL)
	err = file.CreateInFile(mask, originalURL)
	if err != nil {
		sugar.Warnln("Failed write to file %v", err)
	}
	return mask, nil
}

func GetByOriginalURL(s ShortenService, originalURL string) (string, error) {
	result, err := s.Repo.GetByOriginalURL(originalURL)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}
