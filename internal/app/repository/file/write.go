package file

import (
	"encoding/json"
	"fmt"
)

func (handler *Handler) Write(mask string, originalURL string) (string, error) {
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
