package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/maitesin/hermes/config"
	"github.com/maitesin/hermes/internal/app"
	sqlx "github.com/maitesin/hermes/internal/infra/sql"
	"github.com/maitesin/hermes/pkg/comm/telegram"
	"github.com/maitesin/hermes/pkg/tracker/correos"
	"github.com/upper/db/v4/adapter/postgresql"
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

	dbConn, err := sql.Open("postgres", cfg.SQL.DatabaseURL())
	if err != nil {
		log.Panic(err)
		return
	}
	defer dbConn.Close()

	pgConn, err := postgresql.New(dbConn)
	if err != nil {
		log.Panic(err)
		return
	}
	defer pgConn.Close()

	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Panic(err)
		return
	}

	migrations, err := migrate.NewWithSourceInstance("iofs", d, cfg.SQL.DatabaseURL())
	if err != nil {
		log.Panic(err)
		return
	}

	err = migrations.Up()
	if err != nil && err.Error() != "no change" {
		log.Panic(err)
		return
	}

	deliveriesRepository := sqlx.NewDeliveriesRepository(pgConn)

	httpClient := http.DefaultClient
	correosTracker, err := correos.NewTracker(httpClient)
	if err != nil {
		log.Panic(err)
		return
	}

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
		err := listener.Listen(ctx, app.Listen(ctx, correosTracker, deliveriesRepository))
		if err != nil {
			log.Println(err)
			return
		}
	}()

	ticker := time.NewTicker(time.Minute)
	messenger, err := telegram.NewMessenger(ctx, telegramBot)
	if err != nil {
		log.Panic(err)
		return
	}

	err = app.Checker(ctx, ticker, correosTracker, deliveriesRepository, messenger)
	if err != nil {
		log.Panic(err)
		return
	}

	<-ctx.Done()
}
