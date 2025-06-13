package file

import (
	"bufio"
	"os"
)

// URLStorageFileData contains uuid, short URL and original URL
type URLStorageFileData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type handler struct {
	file   *os.File
	writer *bufio.Writer
	reader *bufio.Reader
	Count  int
}

// DBHandler contains info for write and read file
var DBHandler = handler{}

// CreateInFile create record in file
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
