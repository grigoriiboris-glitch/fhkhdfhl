package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymindmap/api/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config, log *log.Logger) *pgxpool.Pool {
	dbpool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}
	return dbpool
}
