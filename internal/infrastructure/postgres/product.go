package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Lovodia/Product/internal/domain"
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
		slog.Error("query product failed", domain.LogErr(err))
		return nil, fmt.Errorf("product repo GetAll: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID, &p.CreatedAt); err != nil {
			slog.Error("scan product failed", domain.LogErr(err))
			return nil, fmt.Errorf("product repo Scan: %w", err)
		}
		products = append(products, p)
	}
	slog.Info("product fetched", slog.Int("count", len(products)))
	return products, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id int) (domain.Product, error) {
	var p domain.Product
	err := r.db.QueryRow(ctx,
		"SELECT id, name, price, category_id, created_at FROM products WHERE id = $1", id).Scan(
		&p.ID, &p.Name, &p.Price, &p.CategoryID, &p.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("product not found", slog.Int("id", id))
			return domain.Product{}, domain.ErrNotFound
		}
		slog.Error("query product failed", domain.LogErr(err), slog.Int("id", id))
		return domain.Product{}, fmt.Errorf("product repo GetByID: %w", err)
	}
	slog.Info("product fetched", slog.Int("id", id))
	return p, nil
}

func (r *ProductRepo) Create(ctx context.Context, p domain.Product) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		"INSERT INTO products(name, price, category_id, created_at) VALUES($1, $2, $3, NOW()) RETURNING id, created_at",
		p.Name, p.Price, p.CategoryID).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		slog.Error("failed to create product", domain.LogErr(err))
		return 0, fmt.Errorf("product repo Create: %w", err)
	}
	slog.Info("product created", slog.Int("id", id), slog.String("name", p.Name))
	return id, nil
}

func (r *ProductRepo) Update(ctx context.Context, id int, p domain.Product) (bool, error) {
	cmdTag, err := r.db.Exec(ctx,
		"UPDATE products SET name = $1, price = $2, category_id = $3 WHERE id = $4",
		p.Name, p.Price, p.CategoryID, id,
	)
	if err != nil {
		slog.Error("failed to update product", domain.LogErr(err))
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
		slog.Error("failed to delete product", domain.LogErr(err))
		return false, fmt.Errorf("delete product id=%d failed: %w", id, err)
	}
	if cmdTag.RowsAffected() == 0 {
		slog.Warn("product not found to delete", slog.Int("id", id))
		return false, domain.ErrNotFound
	}
	slog.Info("product deleted", slog.Int("id", id))
	return true, nil
}
