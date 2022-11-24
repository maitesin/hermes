package telegram

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maitesin/hermes/pkg/comm"
)

type Listener struct {
	ctx context.Context
	api *tgbotapi.BotAPI
}

func NewListener(ctx context.Context, api *tgbotapi.BotAPI) (*Listener, error) {
	return &Listener{
		ctx: ctx,
		api: api,
	}, nil
}

func (l *Listener) Listen(ctx context.Context, handler comm.Handler) error {
	updatesChannel := l.getUpdatesChannel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updatesChannel:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			log.Printf("[%s](%d) %q", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)
			msg := comm.Message{Conversation: update.Message.Chat.ID, Text: update.Message.Text}
			err := handler(msg)
			if err != nil {
				log.Printf("handler failed with %v for message %#v", err, msg)
			}
		}
	}
}

func (l *Listener) getUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return l.api.GetUpdatesChan(u)
}
