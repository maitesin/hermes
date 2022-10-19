package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/maitesin/hermes/config"
	"github.com/maitesin/hermes/internal/app"
	"github.com/maitesin/hermes/internal/infra/migrations"
	sqlx "github.com/maitesin/hermes/internal/infra/sql"
	"github.com/maitesin/hermes/pkg/comm/telegram"
	"github.com/maitesin/hermes/pkg/tracker/correos"
	"github.com/upper/db/v4/adapter/postgresql"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
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

	dbDriver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Panic(err)
		return
	}

	migrations.RegisterMigrationDriver(migrationsFS)
	migrations, err := migrate.NewWithDatabaseInstance("embed://migrations", "marvin", dbDriver)
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

	listener, err := telegram.NewListener(ctx, cfg.Telegram)
	if err != nil {
		log.Panic(err)
		return
	}

	err = listener.Listen(app.Listen(ctx, correosTracker, deliveriesRepository))
	if err != nil {
		log.Panic(err)
		return
	}
}
