package db

import (
	"database/sql"
	"fmt"

	"github.com/Lovodia/Product/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(cfg *config.Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}
	return nil
}
