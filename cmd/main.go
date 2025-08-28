package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Lovodia/Product/internal/config"
	"github.com/Lovodia/Product/internal/db"
	httpdelivery "github.com/Lovodia/Product/internal/delivery/http"
	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/infrastructure"
	"github.com/Lovodia/Product/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	os.Exit(run())
}

func run() int {
	const (
		shutdownTimeout   = 10 * time.Second
		readHeaderTimeout = 5 * time.Second
	)

	migrate := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to Load config:", domain.LogErr(err))

		return 1
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.LogLevel))); err != nil {
		slog.Warn("invalid log level, falling back to info", slog.String("provited", cfg.LogLevel))
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
	slog.Info("Configuration loaded", slog.Any("config", cfg))

	ctx := context.Background()

	if *migrate {
		slog.Info("Running migrations...")

		pool, err := pgxpool.New(ctx, cfg.DSN())

		if err != nil {
			slog.Error("failed to connect for migrations", domain.LogErr(err))

			return 1
		}

		defer pool.Close()

		if err := pool.Ping(ctx); err != nil {
			slog.Error("database ping failed", domain.LogErr(err))

			return 1
		}

		slog.Info("Database ping successful")

		if err := db.RunMigrations(pool, cfg.MigrationsPath); err != nil {
			slog.Error("migrations failed", domain.LogErr(err))

			return 1
		}

		slog.Info("Migrations applied successfully")

		return 0
	}

	database, err := db.New(cfg)
	if err != nil {
		slog.Error("failed to connect to BD", domain.LogErr(err))

		return 1
	}

	defer database.Close()
	slog.Info("Database connection established")

	infrFactory := infrastructure.NewRepositoryFactory(database.Pool())
	productUC := usecase.NewProductUseCase(infrFactory.ProductRepo)
	categoryUC := usecase.NewCategoryUseCase(infrFactory.CategoryRepo)

	r := mux.NewRouter()
	httpdelivery.NewProductHandler(r, productUC)
	httpdelivery.NewCategoryHandler(r, categoryUC)

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	errChan := make(chan error, 1)

	go func() {
		slog.Info("Server starting", slog.String("addr", srv.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		slog.Warn("Shutting down server...", slog.String("signal", sig.String()))
	case err := <-errChan:
		slog.Error("server failed", domain.LogErr(err))

		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", domain.LogErr(err))

		return 1
	}

	slog.Info("Server exited gracefully")

	return 0
}
