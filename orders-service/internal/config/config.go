package config

import (
	"os"
	"strconv"
	"time"

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

	JWTSecret string

	OutboxPollInterval time.Duration
	OutboxBatchSize    int
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

		JWTSecret: getEnv("JWT_SECRET", "dev_secret"),

		OutboxPollInterval: time.Duration(getEnvInt("OUTBOX_POLL_INTERVAL_MS", 500)) * time.Millisecond,
		OutboxBatchSize:    getEnvInt("OUTBOX_BATCH_SIZE", 50),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return defaultValue
	}

	return n
}
