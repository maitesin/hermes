package config

import (
	"os"

	"github.com/maitesin/hermes/pkg/comm/telegram"
)

type Config struct {
	Telegram telegram.Config
}

func NewConfig() (Config, error) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		return Config{}, Error{reason: "TELEGRAM_TOKEN not set"}
	}

	return Config{
		Telegram: telegram.Config{
			Token: telegramToken,
		},
	}, nil
}
