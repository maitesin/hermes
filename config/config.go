package config

import "github.com/maitesin/hermes/pkg/comm/telegram"

type Config struct {
	Telegram telegram.Config
}

func NewConfig() (Config, error) {
	return Config{}, nil
}
