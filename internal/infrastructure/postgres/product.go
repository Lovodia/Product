package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) GetAll(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, price, category_id, created_at FROM products")
	if err != nil {
		slog.Error("query product failed", logger.LogErr(err))

		return nil, fmt.Errorf("product repo GetAll: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0)

	for rows.Next() {
		var prod domain.Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Price, &prod.CategoryID, &prod.CreatedAt); err != nil {
			slog.Error("scan product failed", logger.LogErr(err))

			return nil, fmt.Errorf("product repo Scan: %w", err)
		}

		products = append(products, prod)
	}

	slog.Info("product fetched", slog.Int("count", len(products)))

	return products, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id int) (domain.Product, error) {
	var prod domain.Product
	err := r.db.QueryRow(ctx,
		"SELECT id, name, price, category_id, created_at FROM products WHERE id = $1", id).Scan(
		&prod.ID, &prod.Name, &prod.Price, &prod.CategoryID, &prod.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("product not found", slog.Int("id", id))

			return domain.Product{}, domain.ErrNotFound
		}

		slog.Error("query product failed", logger.LogErr(err), slog.Int("id", id))

		return domain.Product{}, fmt.Errorf("product repo GetByID: %w", err)
	}

	slog.Info("product fetched", slog.Int("id", id))

	return prod, nil
}

func (r *ProductRepo) Create(ctx context.Context, prod domain.Product) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		"INSERT INTO products(name, price, category_id, created_at) VALUES($1, $2, $3, NOW()) RETURNING id, created_at",
		prod.Name, prod.Price, prod.CategoryID).Scan(&id, &prod.CreatedAt)

	if err != nil {
		slog.Error("failed to create product", logger.LogErr(err))

		return 0, fmt.Errorf("product repo Create: %w", err)
	}

	slog.Info("product created", slog.Int("id", id), slog.String("name", prod.Name))

	return id, nil
}

func (r *ProductRepo) Update(ctx context.Context, id int, prod domain.Product) (bool, error) {
	cmdTag, err := r.db.Exec(ctx,
		"UPDATE products SET name = $1, price = $2, category_id = $3 WHERE id = $4",
		prod.Name, prod.Price, prod.CategoryID, id,
	)
	if err != nil {
		slog.Error("failed to update product", logger.LogErr(err))

		return false, fmt.Errorf("update product id=%d failed: %w", id, err)
	}

	if cmdTag.RowsAffected() == 0 {
		slog.Warn("product not found to update", slog.Int("id", id))

		return false, domain.ErrNotFound
	}

	slog.Info("product updated", slog.Int("id", id))

	return true, nil
}

func (r *ProductRepo) Delete(ctx context.Context, id int) (bool, error) {
	cmdTag, err := r.db.Exec(ctx, "DELETE FROM products WHERE id =$1", id)
	if err != nil {
		slog.Error("failed to delete product", logger.LogErr(err))

		return false, fmt.Errorf("delete product id=%d failed: %w", id, err)
	}

	if cmdTag.RowsAffected() == 0 {
		slog.Warn("product not found to delete", slog.Int("id", id))

		return false, domain.ErrNotFound
	}

	slog.Info("product deleted", slog.Int("id", id))

	return true, nil
}
