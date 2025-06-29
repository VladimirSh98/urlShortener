package file

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadLine(t *testing.T) {

	t.Run("read records", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_read.json")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		testRecords := []URLStorageFileData{
			{UUID: "1", ShortURL: "abc123", OriginalURL: "https://example.com/1"},
			{UUID: "2", ShortURL: "def456", OriginalURL: "https://example.com/2"},
		}

		for _, record := range testRecords {
			var data []byte
			data, err = json.Marshal(record)
			require.NoError(t, err)
			_, err = tmpFile.Write(append(data, '\n'))
			require.NoError(t, err)
		}

		tmpFile.Close()
		var file *os.File
		file, err = os.Open(tmpFile.Name())
		require.NoError(t, err)
		defer file.Close()

		handlerTest := &handler{
			reader: bufio.NewReader(file),
		}

		var record1 *URLStorageFileData
		record1, err = handlerTest.ReadLine()
		require.NoError(t, err)
		assert.Equal(t, &testRecords[0], record1)

		var record2 *URLStorageFileData
		record2, err = handlerTest.ReadLine()
		require.NoError(t, err)
		assert.Equal(t, &testRecords[1], record2)

		var record3 *URLStorageFileData
		record3, err = handlerTest.ReadLine()
		require.NoError(t, err)
		assert.Nil(t, record3)
	})
}
