package repository

import "go.uber.org/zap"

var globalURLStorage = map[string]string{}

func Create(mask string, originalURL string, writeToFile bool) string {
	globalURLStorage[mask] = originalURL
	if writeToFile {
		sugar := zap.S()
		err := DBHandler.Open()
		defer DBHandler.Close()
		if err != nil {
			sugar.Warnln("Failed to open file")
			return mask
		}
		_, err = DBHandler.Write(mask, originalURL)
		if err != nil {
			sugar.Warnln("Failed write to file")
			return mask
		}
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
