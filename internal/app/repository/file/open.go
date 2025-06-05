package file

import (
	"bufio"
	"os"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
)

func (handler *Handler) Open() error {
	file, err := os.OpenFile(config.DBFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	handler.file = file
	handler.writer = bufio.NewWriter(handler.file)
	return nil
}
