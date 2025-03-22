package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
)

var DBHandler = FileHandler{}

func (handler *FileHandler) Write(mask string, originalURL string) (string, error) {
	num := fmt.Sprintf("%d", handler.Count+1)
	data := URLStorageFileData{num, mask, originalURL}
	jsonData, err := json.Marshal(data)
	jsonData = append(jsonData, '\n')
	if err != nil {
		return "", err
	}
	_, err = handler.writer.Write(jsonData)
	if err != nil {
		return "", err
	}
	handler.Count++
	err = handler.writer.Flush()
	if err != nil {
		return "", err
	}
	return mask, nil
}

func (handler *FileHandler) Close() error {
	err := handler.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (handler *FileHandler) Open() error {
	path := filepath.Join(config.DBFilePath, config.DBFileName)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	handler.file = file
	handler.writer = bufio.NewWriter(handler.file)
	sugar := zap.S()
	sugar.Infow(path)
	return nil
}

func (handler *FileHandler) OpenReadOnly() error {
	var err error
	_, err = os.Stat(config.DBFilePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(config.DBFilePath, os.ModePerm)
		if err != nil {
			fmt.Printf("Ошибка при создании директории: %v \n", err)
			return err
		}
	}
	fullPath := filepath.Join(config.DBFilePath, config.DBFileName)
	handler.file, err = os.OpenFile(fullPath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	handler.reader = bufio.NewReader(handler.file)
	return nil
}

func (handler *FileHandler) ReadLine() (*URLStorageFileData, error) {
	data, err := handler.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}
	record := URLStorageFileData{}
	err = json.Unmarshal(data, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}
