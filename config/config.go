package config

import (
	"os"

	"github.com/maitesin/hermes/internal/infra/sql"

	"github.com/maitesin/hermes/pkg/comm/telegram"
)

type Config struct {
	Telegram  telegram.Config
	SQL       sql.Config
	DHLAPIKey string
}

func New() (Config, error) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		return Config{}, Error{reason: "TELEGRAM_TOKEN not set"}
	}

	dhlAPIKey := os.Getenv("DHL_API_KEY")
	if telegramToken == "" {
		return Config{}, Error{reason: "DHL_API_KEY not set"}
	}

	return Config{
		Telegram: telegram.Config{
			Token: telegramToken,
		},
		SQL: sql.Config{
			URL:          getEnvOrDefault("DB_URL", "postgres://postgres:postgres@localhost:54321/hermes"),
			SSLMode:      getEnvOrDefault("DB_SSL_MODE", "disable"),
			BinaryParams: getEnvOrDefault("DB_BINARY_PARAMETERS", "yes"),
		},
		DHLAPIKey: dhlAPIKey,
	}, nil
}

func getEnvOrDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}

	return defaultValue
}
