package file

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenReadOnly(t *testing.T) {

	t.Run("successful open existing file", func(t *testing.T) {
		testData := `{"uuid":"1","short_url":"ETe5ORyc","original_url":"http://test.com"}` + "\n"

		filePath := "test_data.json"

		err := os.WriteFile(filePath, []byte(testData), 0644)
		require.NoError(t, err)
		defer os.Remove(filePath)
		require.NoError(t, err)

		config.DBFilePath = filePath

		handlerTest := &handler{}
		err = handlerTest.OpenReadOnly()
		require.NoError(t, err)
		defer handlerTest.file.Close()

		assert.NotNil(t, handlerTest.file)
		assert.NotNil(t, handlerTest.reader)
		assert.FileExists(t, filePath)
	})
}
