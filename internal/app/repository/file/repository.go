package file

import (
	"bufio"
	"os"
)

type URLStorageFileData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Handler struct {
	file   *os.File
	writer *bufio.Writer
	reader *bufio.Reader
	Count  int
}

var DBHandler = Handler{}

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
