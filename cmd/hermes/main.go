package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	
	"github.com/maitesin/hermes/config"
	"github.com/maitesin/hermes/internal/app"
	"github.com/maitesin/hermes/internal/infra/httpx"
	sqlx "github.com/maitesin/hermes/internal/infra/sql"
	"github.com/maitesin/hermes/pkg/comm/telegram"
	"github.com/maitesin/hermes/pkg/tracker/correos"
	"github.com/maitesin/hermes/pkg/tracker/dhl"
	"github.com/maitesin/hermes/pkg/tracker/group"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Panic(err)
		return
	}

	dMigrations, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		panic(err)
	}

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		panic(err)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		panic(err)
	}

	dbDriver, err := migrate.NewWithInstance("iofs", dMigrations, "", driver)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		panic(err)
	}

	err = dbDriver.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Panic(err)
		return
	}

	deliveriesRepository := sqlx.NewDeliveriesRepository(db)

	httpClient := http.DefaultClient
	correosTracker, err := correos.NewTracker(httpClient)
	if err != nil {
		log.Panic(err)
		return
	}

	dhlTracker, err := dhl.NewTracker(httpClient, cfg.DHLAPIKey)
	if err != nil {
		log.Panic(err)
		return
	}

	groupTracker := group.NewGroup(correosTracker, dhlTracker)

	telegramBot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
		return
	}

	listener, err := telegram.NewListener(ctx, telegramBot)
	if err != nil {
		log.Panic(err)
		return
	}

	go func() {
		err := listener.Listen(ctx, app.Listen(ctx, groupTracker, deliveriesRepository))
		if err != nil {
			log.Println(err)
			return
		}
	}()

	go func() {
		err := http.ListenAndServe("0.0.0.0:8181", httpx.ListUndelivered(deliveriesRepository))
		if err != nil {
			log.Println(err)
			return
		}
	}()

	ticker := time.NewTicker(10 * time.Minute)
	messenger, err := telegram.NewMessenger(ctx, telegramBot)
	if err != nil {
		log.Panic(err)
		return
	}

	err = app.Checker(ctx, ticker, groupTracker, deliveriesRepository, messenger)
	if err != nil {
		log.Panic(err)
		return
	}

	<-ctx.Done()
}
