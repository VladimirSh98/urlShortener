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
	ID          string
	OriginalURL string
	CreatedAt   time.Time
	UserID      int
}

type ShortenBatchRequest struct {
	URL    string
	Mask   string
	UserID int
}
