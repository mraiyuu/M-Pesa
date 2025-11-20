package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/mraiyuu/M-Pesa/internal/env"
)

func main() {
	ctx := context.Background()

	/* load env first */
	if err := godotenv.Load(); err != nil {
		slog.Error("warning: failed to load env")
	}

	cfg := config{
		addr: ":8000",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=mpesa sslmode=disable"),
		},
	}

	api := application{
		config: cfg,
	}

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	logger.Info("connected to database")

	if err := api.run(api.mount()); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
