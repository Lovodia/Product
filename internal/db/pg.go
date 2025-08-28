package db

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/Lovodia/Product/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Config) (*Database, error) {
	ctx := context.Background()

	hostPort := net.JoinHostPort(cfg.DBHost, strconv.Itoa(cfg.DBPort))

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, hostPort, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{pool: pool}, nil
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *Database) Close() {
	db.pool.Close()
}
