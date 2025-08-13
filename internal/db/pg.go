package db

import (
	"context"
	"fmt"

	"github.com/Lovodia/Product/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}
	return &Database{pool: pool}, nil
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *Database) Close() {
	db.pool.Close()
}
