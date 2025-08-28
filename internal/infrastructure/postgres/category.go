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

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetAll(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name FROM categories")
	if err != nil {
		slog.Error("query categories failed", domain.LogErr(err))

		return nil, fmt.Errorf("query categories failed: %w", err)
	}
	defer rows.Close()

	categories := make([]domain.Category, 0)

	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			slog.Error("scan category failed", domain.LogErr(err))

			return nil, fmt.Errorf("scan category failed: %w", err)
		}

		categories = append(categories, cat)
	}

	slog.Info("categories fetched", slog.Int("count", len(categories)))

	return categories, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int) (domain.Category, error) {
	var cat domain.Category
	err := r.db.QueryRow(ctx, "SELECT id, name FROM categories WHERE id=$1", id).
		Scan(&cat.ID, &cat.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("category not found", slog.Int("id", id))

			return domain.Category{}, domain.ErrNotFound
		}

		slog.Error("query category failed", domain.LogErr(err), slog.Int("id", id))

		return domain.Category{}, fmt.Errorf("category repo GetByID: %w", err)
	}

	slog.Info("category fetched", slog.Int("id", id))

	return cat, nil
}

func (r *CategoryRepo) Create(ctx context.Context, category domain.Category) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		"INSERT INTO categories(name) VALUES($1) RETURNING id", category.Name).Scan(&id)

	if err != nil {
		slog.Error("failed to create category", domain.LogErr(err))

		return 0, fmt.Errorf("category repo Create: %w", err)
	}

	slog.Info("category created", slog.Int("id", id), slog.String("name", category.Name))

	return id, nil
}

func (r *CategoryRepo) Update(ctx context.Context, category domain.Category) (bool, error) {
	cmdTag, err := r.db.Exec(ctx,
		"UPDATE categories SET name=$1 WHERE id=$2", category.Name, category.ID)
	if err != nil {
		slog.Error("failed to update category", domain.LogErr(err))

		return false, fmt.Errorf("update category id=%d failed: %w", category.ID, err)
	}

	if cmdTag.RowsAffected() == 0 {
		slog.Warn("category not found to update", slog.Int("id", category.ID))

		return false, domain.ErrNotFound
	}

	slog.Info("category updated", slog.Int("id", category.ID))

	return true, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id int) (bool, error) {
	cmdTag, err := r.db.Exec(ctx, "DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		slog.Error("failed to delete category", domain.LogErr(err))

		return false, fmt.Errorf("delete category id=%d failed: %w", id, err)
	}

	if cmdTag.RowsAffected() == 0 {
		slog.Warn("category not found to delete", slog.Int("id", id))

		return false, domain.ErrNotFound
	}

	slog.Info("category deleted", slog.Int("id", id))

	return true, nil
}
