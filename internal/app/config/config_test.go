package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {

	t.Run("env variables", func(t *testing.T) {
		os.Setenv("SERVER_ADDRESS", "127.0.0.1:8081")
		os.Setenv("BASE_URL", "http://127.0.0.1:8081")
		os.Setenv("DB_FILE_PATH", "./db.json")
		os.Setenv("DATABASE_DSN", "postgres://env:env@localhost:5432/env")

		err := LoadConfig()
		require.NoError(t, err)

		assert.Equal(t, "127.0.0.1:8081", FlagRunAddr)
		assert.Equal(t, "http://127.0.0.1:8081", FlagResultAddr)
		assert.Equal(t, "./db.json", DBFilePath)
		assert.Equal(t, "postgres://env:env@localhost:5432/env", DatabaseDSN)

		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("BASE_URL")
		os.Unsetenv("DB_FILE_PATH")
		os.Unsetenv("DATABASE_DSN")
	})

	t.Run("command line flags", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{"cmd", "-a", "localhost:8080", "-b", "http://localhost:8080", "-f", "/tmp/flag_db.json", "-d", "postgres://flag:flag@localhost:5432/flag"}

		err := LoadConfig()
		require.NoError(t, err)

		assert.Equal(t, "localhost:8080", FlagRunAddr)
		assert.Equal(t, "http://localhost:8080", FlagResultAddr)
		assert.Equal(t, "/tmp/flag_db.json", DBFilePath)
		assert.Equal(t, "postgres://flag:flag@localhost:5432/flag", DatabaseDSN)
	})

	t.Run("default yaml config", func(t *testing.T) {
		origFile := "default_config.yaml"
		var err error
		var origData []byte
		origData, err = os.ReadFile(origFile)
		if err == nil {
			defer os.WriteFile(origFile, origData, 0644)
		}

		require.NoError(t, err)

		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{"cmd"}

		err = LoadConfig()
		require.NoError(t, err)

		assert.Equal(t, "localhost:8080", FlagRunAddr)
		assert.Equal(t, "http://localhost:8080", FlagResultAddr)
		assert.Equal(t, "./db.json", DBFilePath)
		assert.Equal(t, "postgres://user:zhoskiy@localhost:5432/shortner?sslmode=disable", DatabaseDSN)
	})
}
