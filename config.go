package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr      string
	SecretKey []byte
	Database  *DatabaseConfig
}

func New(opts ...func(*Config)) *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found: %v", err)
	}

	cfg := &Config{
		Addr:      GetEnv("ADDR", ":8080"),
		SecretKey: []byte(GetEnv("SECRET_KEY", "")),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func WithDatabaseConfig(cfg *DatabaseConfig) func(c *Config) {
	return func(c *Config) {
		c.Database = cfg
	}
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("environment variable %s is not set", key)
		return defaultValue
	}

	return value

}

func GetEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("environment variable %s is not set, using default %d", key, defaultValue)
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("error converting environment variable %s to int: %v", key, err)
		return defaultValue
	}
	return intValue
}
