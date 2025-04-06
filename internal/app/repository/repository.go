package repository

import "go.uber.org/zap"

var globalURLStorage = map[string]string{}

func Create(mask string, originalURL string) string {
	var err error
	sugar := zap.S()
	globalURLStorage[mask] = originalURL
	_, err = createDB(mask, originalURL)
	if err != nil {
		sugar.Warnln("Failed write to database %v", err)
	}
	err = CreateInFile(mask, originalURL)
	if err != nil {
		sugar.Warnln("Failed write to file %v", err)
	}
	return mask
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
