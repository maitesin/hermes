package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maitesin/hermes/pkg/messenger"
)

type Messenger struct {
	config Config
	ctx    context.Context
	api    *tgbotapi.BotAPI
}

func NewMessenger(ctx context.Context, cfg Config) (*Messenger, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		fmt.Println(err)
		return &Messenger{}, err
	}

	api.Debug = true
	return &Messenger{
		ctx:    ctx,
		api:    api,
		config: cfg,
	}, nil
}

func (m *Messenger) Message(message messenger.Message) error {
	_, err := m.api.Send(tgbotapi.NewMessage(message.Conversation, message.Text))
	return err
}
