package config


func New() (*Config, error) {
	config := Config{}

	config.DatabaseName = GetEnv("DB_NAME", "unitrip")
	config.DatabaseUser = GetEnv("DB_USER", "")
	config.DatabasePassword = GetEnv("DB_PASSWORD", "")
	config.DatabasePort = GetEnv("DB_PORT", "5432")
	config.DatabaseHost = GetEnv("DB_HOST", "localhost")

	config.WebPort = GetEnv("WEB_PORT", "8080")
	config.SecretKey = GetEnv("SECRET_KEY", "")
	return &config, nil
}
