package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Lovodia/Product/internal/config"
	"github.com/Lovodia/Product/internal/db"
	httpDelivery "github.com/Lovodia/Product/internal/delivery/http"
	"github.com/Lovodia/Product/internal/infrastructure"
	"github.com/Lovodia/Product/internal/usecase"
	"github.com/gorilla/mux"
)

func main() {
	migrate := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to Load config:", slog.Any("err", err))
		os.Exit(1)
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.LogLevel))); err != nil {
		slog.Warn("invalid log level, falling back to info", slog.String("provited", cfg.LogLevel))
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	slog.Info("Configuration loaded", slog.Any("config", cfg))

	if *migrate {
		slog.Info("Running migrations...")
		if err := db.RunMigrations(cfg); err != nil {
			slog.Error("migrations failed", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Migrations applied successfully")
		return
	}

	database, err := db.New(cfg)
	if err != nil {
		slog.Error("failed to connect to BD", slog.Any("error", err))
		os.Exit(1)
	}
	defer database.Close()

	slog.Info("Database connection established")

	infrFactory := infrastructure.NewRepositoryFactory(database.Pool())

	productUC := usecase.NewProductUseCase(infrFactory.ProductRepo)
	categoryUC := usecase.NewCategoryUseCase(infrFactory.CategoryRepo)

	productHandler := httpDelivery.NewProductHandler(productUC)
	categoryHandler := httpDelivery.NewCategoryHandler(categoryUC)

	r := mux.NewRouter()

	r.HandleFunc("/product", productHandler.GetAll).Methods("GET")
	r.HandleFunc("/product/{id:[0-9]+}", productHandler.GetByID).Methods("GET")
	r.HandleFunc("/product", productHandler.Create).Methods("POST")
	r.HandleFunc("/product/{id:[0-9]+}", productHandler.Update).Methods("PUT")
	r.HandleFunc("/product/{id:[0-9]+}", productHandler.Delete).Methods("DELETE")

	r.HandleFunc("/category", categoryHandler.GetAll).Methods("GET")
	r.HandleFunc("/category/{id:[0-9]+}", categoryHandler.GetByID).Methods("GET")
	r.HandleFunc("/category", categoryHandler.Create).Methods("POST")
	r.HandleFunc("/category/{id:[0-9]+}", categoryHandler.Update).Methods("PUT")
	r.HandleFunc("/category/{id:[0-9]+}", categoryHandler.Delete).Methods("DELETE")

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		slog.Info("Server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Warn("Sutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Server exited gracefuly")
}
