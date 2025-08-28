package db

import (
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(pool *pgxpool.Pool, migrationsPatch string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Warn("failed to close sqlBD after migrations", slog.Any("error", err))
		}
	}()

	if err := goose.Up(sqlDB, migrationsPatch); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
