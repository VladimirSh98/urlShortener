package file

import (
	"bufio"
	"os"

	"github.com/VladimirSh98/urlShortener/internal/app/config"
)

// OpenReadOnly open file in readonly mode
func (handler *handler) OpenReadOnly() error {
	var err error
	handler.file, err = os.OpenFile(config.DBFilePath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	handler.reader = bufio.NewReader(handler.file)
	return nil
}
