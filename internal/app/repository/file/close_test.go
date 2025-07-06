package file

import (
	"github.com/VladimirSh98/urlShortener/internal/app/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerClose(t *testing.T) {
	origDBPath := config.DBFilePath
	defer func() { config.DBFilePath = origDBPath }()

	t.Run("successful close", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.json")
		file, err := os.Create(dbPath)
		require.NoError(t, err)

		handlerTest := &handler{
			file: file,
		}

		err = handlerTest.Close()
		require.NoError(t, err)

		_, err = file.WriteString("test")
		assert.Error(t, err, "Write to closed file should fail")
		assert.Contains(t, err.Error(), "closed")
	})
}
