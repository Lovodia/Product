package bootstrap

import (
	"flag"
	"log/slog"

	"github.com/Lovodia/Product/internal/config"
)

func LoadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return nil, err
	}
	return cfg, nil
}

func IsMigrationMode() bool {
	migrate := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()
	return *migrate
}
