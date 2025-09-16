package app

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lovodia/Product/internal/config"
	"github.com/Lovodia/Product/internal/db"
	httpdelivery "github.com/Lovodia/Product/internal/delivery/http"
	"github.com/Lovodia/Product/internal/infrastructure"
	"github.com/Lovodia/Product/internal/logger"
	"github.com/Lovodia/Product/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	cfg *config.Config
	db  *db.Database
}

func New() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger.SetupLogger(cfg)

	dbConn, err := db.New(cfg)
	if err != nil {
		return nil, err
	}

	slog.Info("Application initialized")

	return &Application{
		cfg: cfg,
		db:  dbConn,
	}, nil
}

func (a *Application) Close() {
	if a.db != nil {
		a.db.Close()
	}
}

func (a *Application) IsMigrationMode() bool {
	migrate := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	return *migrate
}

func (a *Application) RunMigrations() int {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, a.cfg.DB.DSN())

	if err != nil {
		slog.Error("failed to connect for migrations", logger.LogErr(err))

		return 1
	}

	defer pool.Close()

	if err := db.RunMigrations(pool, a.cfg.Migration.Path); err != nil {
		slog.Error("migrations failed", logger.LogErr(err))

		return 1
	}

	slog.Info("Migrations applied successfully")

	return 0
}

func (a *Application) StartServer() int {
	const (
		shutdownTimeout   = 10 * time.Second
		readHeaderTimeout = 5 * time.Second
	)

	infrFactory := infrastructure.NewRepositoryFactory(a.db.Pool())
	productUC := usecase.NewProductUseCase(infrFactory.ProductRepo)
	categoryUC := usecase.NewCategoryUseCase(infrFactory.CategoryRepo)

	r := mux.NewRouter()
	httpdelivery.NewProductHandler(r, productUC)
	httpdelivery.NewCategoryHandler(r, categoryUC)

	srv := &http.Server{
		Addr:              ":" + a.cfg.Server.Port,
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
		slog.Error("server failed", logger.LogErr(err))

		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", logger.LogErr(err))

		return 1
	}

	slog.Info("Server exited gracefully")

	return 0
}
