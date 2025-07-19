package config

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseURL       string `env:"BASE_URL" json:"base_url"`
	DBFilePath    string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN   string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS   bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

type defaultConfig struct {
	ServerAddress string `yaml:"server_address"`
	BaseURL       string `yaml:"base_url"`
	DBFilePath    string `yaml:"db_file_path"`
	DBFileName    string `yaml:"db_file_name"`
	DatabaseDSN   string `yaml:"database_dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
	CertFile      string `yaml:"cert_file"`
	KeyFile       string `yaml:"key_file"`
	GrpcAddress   string `yaml:"grpc_address"`
}
