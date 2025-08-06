package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Lovodia/Product/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode)

	if err := runMigrations(dsn, "./migrations"); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &Database{Pool: pool}, nil
}

func runMigrations(dsn string, migrationsDir string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}
	return nil
}
