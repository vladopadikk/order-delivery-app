package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	KafkaHost string
	KafkaPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load("../.env")

	cfg := &Config{
		AppPort: getEnv("APP_PORT", "8080"),

		KafkaHost: getEnv("KAFKA_HOST", "localhost"),
		KafkaPort: getEnv("KAFKA_PORT", "9092"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "dev_secret"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
