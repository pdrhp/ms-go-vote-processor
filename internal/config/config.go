package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
}

type AppConfig struct {
	Environment string
	LogLevel    string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	RunMigrations   bool
	MigrationsPath  string
}

type KafkaConfig struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
	BatchSize     int
	Workers       int
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		App: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_NAME", "votes"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", "5m"),
			RunMigrations:   getEnvBool("DB_RUN_MIGRATIONS", true),
			MigrationsPath:  getEnv("DB_MIGRATIONS_PATH", "./internal/infrastructure/persistence/migrations"),
		},
		Kafka: KafkaConfig{
			Brokers:       getEnvSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:         getEnv("KAFKA_TOPIC", "votos"),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "vote-processor"),
			BatchSize:     getEnvInt("KAFKA_BATCH_SIZE", 1000),
			Workers:       getEnvInt("KAFKA_WORKERS", 5),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return 5 * time.Minute
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
