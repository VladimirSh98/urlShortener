package file

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateInFile(t *testing.T) {

	t.Run("successful create and write", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.json")
		config.DBFilePath = dbPath

		DBHandler = handler{}

		err := CreateInFile("abc123", "https://example.com")
		require.NoError(t, err)

		assert.Equal(t, 1, DBHandler.Count)

		content, err := os.ReadFile(dbPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "abc123")
		assert.Contains(t, string(content), "https://example.com")
	})

	t.Run("create multiple records", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.json")
		config.DBFilePath = dbPath
		DBHandler = handler{}
		err := CreateInFile("first", "https://first.com")
		require.NoError(t, err)

		err = CreateInFile("second", "https://second.com")
		require.NoError(t, err)

		content, err := os.ReadFile(dbPath)
		require.NoError(t, err)
		assert.Equal(t, 2, DBHandler.Count)
		assert.Contains(t, string(content), "first")
		assert.Contains(t, string(content), "second")
	})

	t.Run("file cannot be opened", func(t *testing.T) {
		// Невалидный путь
		config.DBFilePath = "/invalid/path/test.json"
		DBHandler = handler{}

		err := CreateInFile("test", "https://test.com")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "open")
	})
}
