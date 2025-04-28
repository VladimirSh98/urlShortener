package config

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	DBFilePath    string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

type DefaultConfig struct {
	ServerAddress string `yaml:"server_address"`
	BaseURL       string `yaml:"base_url"`
	DBFilePath    string `yaml:"db_file_path"`
	DBFileName    string `yaml:"db_file_name"`
	DatabaseDSN   string `yaml:"database_dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}
