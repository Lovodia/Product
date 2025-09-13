package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Lovodia/Product/internal/config"
)

func LogErr(err error) slog.Attr {
	return slog.Any("error", err)
}

func SetupLogger(cfg *config.Config) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Logger.Level))); err != nil {
		slog.Warn("invalid log level, falling back to info", slog.String("provided", cfg.Logger.Level))
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
	slog.Info("configuration loaded", slog.Any("config", cfg))
}
