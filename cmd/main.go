package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Lovodia/Product/internal/config"
	"github.com/Lovodia/Product/internal/db"
	httpDelivery "github.com/Lovodia/Product/internal/delivery/http"
	"github.com/Lovodia/Product/internal/infrastructure/postgres"
	"github.com/Lovodia/Product/internal/usecase"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to Load config:", err)
	}

	database, err := db.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to BD:", err)
	}
	defer database.Pool.Close()

	productRepo := postgres.NewProductRepo(database.Pool)
	productUC := usecase.NewProductUseCase(productRepo)
	productHandler := httpDelivery.NewProductHandler(productUC)

	categoryRepo := postgres.NewCategoryRepo(database.Pool)
	categoryUC := usecase.NewCategoryUseCase(categoryRepo)
	categoryHandler := httpDelivery.NewCategoryHandler(categoryUC)

	r := mux.NewRouter()

	r.HandleFunc("/product", productHandler.GetAll).Methods("GET")
	r.HandleFunc("/product/{id:[0-9]+}", productHandler.GetByID).Methods("GET")
	r.HandleFunc("/product", productHandler.Create).Methods("Post")
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
		log.Println("Server is running at http://localhost:" + cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Sutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefuly")
}
