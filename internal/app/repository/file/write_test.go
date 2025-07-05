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
	})
}
