package postgres

import (
	"context"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
	DB *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) repository.CategoryRepository {
	return &CategoryRepo{DB: db}
}

func (r *CategoryRepo) GetAll() ([]domain.Category, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepo) GetByID(id int) (domain.Category, error) {
	var c domain.Category
	err := r.DB.QueryRow(context.Background(), "SELECT id, name FROM categories WHERE id=$1", id).
		Scan(&c.ID, &c.Name)
	return c, err
}

func (r *CategoryRepo) Create(category domain.Category) (domain.Category, error) {
	var id int
	err := r.DB.QueryRow(context.Background(),
		"INSERT INTO categories(name) VALUES($1) RETURNING id", category.Name).Scan(&id)
	if err != nil {
		return domain.Category{}, err
	}
	category.ID = id
	return category, nil
}

func (r *CategoryRepo) Update(category domain.Category) error {
	_, err := r.DB.Exec(context.Background(),
		"UPDATE categories SET name=$1 WHERE id=$2", category.Name, category.ID)
	return err
}

func (r *CategoryRepo) Delete(id int) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM categories WHERE id=$1", id)
	return err
}
