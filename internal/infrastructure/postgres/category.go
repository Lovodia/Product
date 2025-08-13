package postgres

import (
	"context"
	"fmt"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetAll() ([]domain.Category, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name FROM categories")
	if err != nil {
		return nil, fmt.Errorf("query categories failed: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("scan category failed: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepo) GetByID(id int) (domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRow(context.Background(), "SELECT id, name FROM categories WHERE id=$1", id).
		Scan(&c.ID, &c.Name)
	if err != nil {
		return domain.Category{}, fmt.Errorf("query category by id=%d failed: %w", id, err)
	}
	return c, nil
}

func (r *CategoryRepo) Create(category domain.Category) (int, error) {
	var id int
	err := r.db.QueryRow(context.Background(),
		"INSERT INTO categories(name) VALUES($1) RETURNING id", category.Name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CategoryRepo) Update(category domain.Category) (bool, error) {
	cmdTag, err := r.db.Exec(context.Background(),
		"UPDATE categories SET name=$1 WHERE id=$2", category.Name, category.ID)
	if err != nil {
		return false, err
	}
	return cmdTag.RowsAffected() > 0, nil
}

func (r *CategoryRepo) Delete(id int) (bool, error) {
	cmdTag, err := r.db.Exec(context.Background(), "DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	return cmdTag.RowsAffected() > 0, nil
}
