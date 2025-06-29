package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	origStorage := GlobalURLStorage
	defer func() { GlobalURLStorage = origStorage }()

	t.Run("existing", func(t *testing.T) {
		GlobalURLStorage = map[string]string{
			"abc123": "https://example.com",
		}
		url, ok := Get("abc123")

		assert.True(t, ok)
		assert.Equal(t, "https://example.com", url)
	})

	t.Run("non-existing", func(t *testing.T) {
		GlobalURLStorage = map[string]string{}

		url, ok := Get("not_exist")

		assert.False(t, ok)
		assert.Empty(t, url)
	})
}

func TestDelete(t *testing.T) {
	origStorage := GlobalURLStorage
	defer func() { GlobalURLStorage = origStorage }()

	t.Run("delete existing", func(t *testing.T) {
		GlobalURLStorage = map[string]string{
			"abc123": "https://example.com",
			"def456": "https://another.com",
		}

		Delete("abc123")

		_, exists := GlobalURLStorage["abc123"]
		assert.False(t, exists)
		assert.Equal(t, 1, len(GlobalURLStorage))
	})

	t.Run("delete non-existing", func(t *testing.T) {
		GlobalURLStorage = map[string]string{
			"abc123": "https://example.com",
		}

		Delete("not_exist")

		assert.Equal(t, 1, len(GlobalURLStorage))
	})
}

func TestCreateInMemory(t *testing.T) {
	origStorage := GlobalURLStorage
	t.Cleanup(func() { GlobalURLStorage = origStorage })

	t.Run("create new", func(t *testing.T) {
		GlobalURLStorage = make(map[string]string)

		CreateInMemory("abc123", "https://example.com")

		assert.Equal(t, 1, len(GlobalURLStorage))
		assert.Equal(t, "https://example.com", GlobalURLStorage["abc123"])
	})

	t.Run("overwrite existing", func(t *testing.T) {
		GlobalURLStorage = map[string]string{
			"abc123": "https://old.example.com",
		}

		CreateInMemory("abc123", "https://new.example.com")

		assert.Equal(t, 1, len(GlobalURLStorage))
		assert.Equal(t, "https://new.example.com", GlobalURLStorage["abc123"])
	})
}
