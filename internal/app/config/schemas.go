package config

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	DBFilePath    string `env:"DB_FILE_PATH"`
}

type defaultConfig struct {
	ServerAddress string `yaml:"server_address"`
	BaseURL       string `yaml:"base_url"`
	DBFilePath    string `yaml:"db_file_path"`
}
