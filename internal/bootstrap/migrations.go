package bootstrap

import (
	"context"
	"log/slog"

	"github.com/Lovodia/Product/internal/config"
	"github.com/Lovodia/Product/internal/db"
	"github.com/Lovodia/Product/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(cfg *config.Config) int {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DB.DSN())
	if err != nil {
		slog.Error("failed to connect for migrations", logger.LogErr(err))
		return 1
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("database ping failed", logger.LogErr(err))
		return 1
	}

	slog.Info("Database ping successful")

	if err := db.RunMigrations(pool, cfg.Migration.Path); err != nil {
		slog.Error("migrations failed", logger.LogErr(err))
		return 1
	}

	slog.Info("Migrations applied successfully")
	return 0
}
