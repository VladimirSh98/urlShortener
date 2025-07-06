package file

import (
	"bufio"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {

	t.Run("successful open file", func(t *testing.T) {
		testData := `{"uuid":"1","short_url":"ETe5ORyc","original_url":"http://test.com"}` + "\n"

		filePath := "test_data.json"

		err := os.WriteFile(filePath, []byte(testData), 0644)
		require.NoError(t, err)
		defer os.Remove(filePath)
		require.NoError(t, err)

		config.DBFilePath = filePath

		handlerTest := &handler{}
		err = handlerTest.Open()
		require.NoError(t, err)
		defer handlerTest.file.Close()

		assert.NotNil(t, handlerTest.file)
		assert.NotNil(t, handlerTest.writer)
		assert.IsType(t, &bufio.Writer{}, handlerTest.writer)
		assert.FileExists(t, filePath)
	})

	t.Run("error on invalid path", func(t *testing.T) {
		config.DBFilePath = "/nonexistent/path/test.json"

		handlerTest := &handler{}
		err := handlerTest.Open()
		require.Error(t, err)
		assert.Nil(t, handlerTest.file)
		assert.Nil(t, handlerTest.writer)
	})
}
