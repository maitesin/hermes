package telegram

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maitesin/hermes/pkg/comm"
)

type Listener struct {
	ctx context.Context
	api *tgbotapi.BotAPI
}

func NewListener(ctx context.Context, cfg Config) (*Listener, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		fmt.Println(err)
		return &Listener{}, err
	}

	api.Debug = true
	return &Listener{
		ctx: ctx,
		api: api,
	}, nil
}

func (l *Listener) Listen(ctx context.Context, handler comm.Handler) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-l.getUpdatesChannel():
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
