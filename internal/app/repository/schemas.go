package repository

import (
	"bufio"
	"os"
)

type URLStorageFileData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileHandler struct {
	file   *os.File
	writer *bufio.Writer
	reader *bufio.Reader
	Count  int
}
