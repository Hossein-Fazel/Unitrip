package config

type Config struct {
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     string
	DatabaseHost     string

	WebPort string

	SecretKey string

	AdminUsername string
	AdminPassword string
}
