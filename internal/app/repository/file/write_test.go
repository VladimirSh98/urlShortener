package file

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_write.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	handlerTest := &handler{
		writer: bufio.NewWriter(tmpFile),
		Count:  0,
	}

	t.Run("successful write", func(t *testing.T) {
		mask := "abc123"
		originalURL := "https://example.com"
		var resultMask string
		resultMask, err = handlerTest.Write(mask, originalURL)
		require.NoError(t, err)

		assert.Equal(t, mask, resultMask)
		assert.Equal(t, 1, handlerTest.Count)

		//// Читаем записанные данные из файла
		//_, err = tmpFile.Seek(0, 0)
		//require.NoError(t, err)
		//
		//var data URLStorageFileData
		//decoder := json.NewDecoder(tmpFile)
		//err = decoder.Decode(&data)
		//require.NoError(t, err)
		//
		//// Проверяем записанные данные
		//assert.Equal(t, "1", data.Num) // handler.Count+1
		//assert.Equal(t, mask, data.Mask)
		//assert.Equal(t, originalURL, data.OriginalURL)
	})
}
