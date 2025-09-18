package main

import (
	"log/slog"
	"os"

	"github.com/Lovodia/Product/internal/app"
)

func main() {
	os.Exit(run())
}

func run() int {
	application, err := app.New()
	if err != nil {
		slog.Error("failed to initialize application", slog.Any("error", err))

		return 1
	}
	defer application.Close()

	if application.IsMigrationMode() {
		return application.RunMigrations()
	}

	return application.StartServer()
}
