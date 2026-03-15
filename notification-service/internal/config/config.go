package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	KafkaHost string
	KafkaPort string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),

		KafkaHost: getEnv("KAFKA_HOST", "localhost"),
		KafkaPort: getEnv("KAFKA_PORT", "9092"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
