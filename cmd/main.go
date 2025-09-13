package main

import (
	"log/slog"
	"os"

	"github.com/Lovodia/Product/internal/bootstrap"
	"github.com/Lovodia/Product/internal/logger"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := bootstrap.LoadConfig()
	if err != nil {
		return 1
	}

	logger.SetupLogger(cfg)

	if bootstrap.IsMigrationMode() {
		return bootstrap.RunMigrations(cfg)
	}

	dbConn, err := bootstrap.InitDatabase(cfg)
	if err != nil {
		slog.Error("failed to connect to DB", logger.LogErr(err))
		return 1
	}
	defer dbConn.Close()

	slog.Info("Database connection established")

	return bootstrap.StartServer(cfg, dbConn)
}
