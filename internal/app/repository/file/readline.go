package file

import (
	"encoding/json"
	"io"
)

func (handler *Handler) ReadLine() (*URLStorageFileData, error) {
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
