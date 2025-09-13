package bootstrap

import (
	"context"
	"errors"
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
)

func StartServer(cfg *config.Config, database *db.Database) int {
	const (
		shutdownTimeout   = 10 * time.Second
		readHeaderTimeout = 5 * time.Second
	)

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
