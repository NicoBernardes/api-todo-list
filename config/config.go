package config

import "os"

// Config DB
type Config struct {
	Server   string
	Port     string
	User     string
	Password string
	database string
}

// LoadConfig
func LoadConfig() *Config {
	return &Config{
		Server:   getEnv("DB_SERVER", "localhost"),
		Port:     getEnv("DB_PORT", "1433"),
		User:     getEnv("DB_USER", "sa"),
		Password: getEnv("DB_PASSWORD", "112233"),
		database: getEnv("DB_NAME", "TodoDB"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
