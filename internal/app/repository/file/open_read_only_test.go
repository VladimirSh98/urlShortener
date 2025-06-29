package file

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenReadOnly(t *testing.T) {
	origDBPath := config.DBFilePath
	defer func() { config.DBFilePath = origDBPath }()

	t.Run("successful open existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.json")
		_, err := os.Create(dbPath)
		require.NoError(t, err)

		config.DBFilePath = dbPath

		handlerTest := &handler{}
		err = handlerTest.OpenReadOnly()
		require.NoError(t, err)
		defer handlerTest.file.Close()

		assert.NotNil(t, handlerTest.file)
		assert.NotNil(t, handlerTest.reader)
		assert.FileExists(t, dbPath)
	})

	t.Run("error on permission denied", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "no_access.db")
		_, err := os.Create(dbPath)
		require.NoError(t, err)
		err = os.Chmod(dbPath, 0222)
		require.NoError(t, err)

		config.DBFilePath = dbPath

		handlerTest := &handler{}
		err = handlerTest.OpenReadOnly()
		require.Error(t, err)
		assert.Nil(t, handlerTest.file)
		assert.Nil(t, handlerTest.reader)
	})
}
