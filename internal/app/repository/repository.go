package repository

import (
	"errors"
	customErr "github.com/VladimirSh98/urlShortener/internal/app/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

var globalURLStorage = map[string]string{}

func Create(mask string, originalURL string, userID int) (string, error) {
	var err error
	sugar := zap.S()
	_, err = createDB(mask, originalURL, userID)
	if err != nil {
		sugar.Warnln("Failed write to database %v", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			sugar.Infoln("URL already exists %s", originalURL)
			var oldMask string
			oldMask, err = GetByOriginalURL(originalURL)
			if err != nil {
				return "", err
			}
			return oldMask, customErr.ErrConstraintViolation
		}
	}
	globalURLStorage[mask] = originalURL
	err = CreateInFile(mask, originalURL)
	if err != nil {
		sugar.Warnln("Failed write to file %v", err)
	}
	return mask, nil
}

func Get(mask string) (string, bool) {
	result, ok := globalURLStorage[mask]
	return result, ok
}

func Delete(mask string) {
	delete(globalURLStorage, mask)
}

func CreateInMemory(mask string, originalURL string) {
	globalURLStorage[mask] = originalURL
}

func CreateInFile(mask string, originalURL string) error {
	err := DBHandler.Open()
	defer DBHandler.Close()
	if err != nil {

		return err
	}
	_, err = DBHandler.Write(mask, originalURL)
	if err != nil {
		return err
	}
	return nil
}

func BatchCreate(data []ShortenBatchRequest) {
	sugar := zap.S()
	err := BatchCreateDB(data)
	if err != nil {
		sugar.Warnln("Failed write to database %v", err)
	}
	for _, record := range data {
		CreateInMemory(record.Mask, record.URL)
		err = CreateInFile(record.Mask, record.URL)
		if err != nil {
			sugar.Warnln("Failed write to file %v", err)
		}
	}
}

func GetByOriginalURL(originalURL string) (string, error) {
	result, err := GetByOriginalURLFromBD(originalURL)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}
