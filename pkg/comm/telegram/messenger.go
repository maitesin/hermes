package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maitesin/hermes/pkg/comm"
)

type Messenger struct {
	ctx context.Context
	api *tgbotapi.BotAPI
}

func NewMessenger(ctx context.Context, api *tgbotapi.BotAPI) (*Messenger, error) {
	return &Messenger{
		ctx: ctx,
		api: api,
	}, nil
}

func (m *Messenger) Message(message comm.Message) error {
	_, err := m.api.Send(tgbotapi.NewMessage(message.Conversation, message.Text))
	return err
}
