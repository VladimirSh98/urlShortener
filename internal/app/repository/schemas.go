package repository

import (
	"bufio"
	"os"
	"time"
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

type Shortner struct {
	Id          string
	OriginalURL string
	CreatedAt   time.Time
}
