package file

import (
	"bufio"
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	origDBPath := config.DBFilePath
	defer func() { config.DBFilePath = origDBPath }()

	t.Run("successful open file", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.json")
		_, err := os.Create(dbPath)
		require.NoError(t, err)

		config.DBFilePath = dbPath

		handlerTest := &handler{}
		err = handlerTest.Open()
		require.NoError(t, err)
		defer handlerTest.file.Close()

		assert.NotNil(t, handlerTest.file)
		assert.NotNil(t, handlerTest.writer)
		assert.IsType(t, &bufio.Writer{}, handlerTest.writer)
		assert.FileExists(t, dbPath)
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
