package main

import (
	"context"
	"log"
	"net/http"

	"github.com/maitesin/hermes/config"
	"github.com/maitesin/hermes/internal/app"
	"github.com/maitesin/hermes/pkg/comm/telegram"
	"github.com/maitesin/hermes/pkg/tracker/correos"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic(err)
		return
	}

	httpClient := http.DefaultClient
	correosTracker, err := correos.NewTracker(httpClient)
	if err != nil {
		log.Panic(err)
		return
	}

	listener, err := telegram.NewListener(ctx, cfg.Telegram)
	if err != nil {
		log.Panic(err)
		return
	}

	err = listener.Listen(app.Listen(correosTracker))
	if err != nil {
		log.Panic(err)
		return
	}
}
