package postgres

import (
	"context"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	DB *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (r *ProductRepo) GetAll() ([]domain.Product, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, name, price, category_id, created_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CategoryID, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepo) GetByID(id int) (domain.Product, error) {
	var p domain.Product
	err := r.DB.QueryRow(context.Background(),
		"SELECT id, name, price, category_id, created_at FROM products WHERE id = $1", id).Scan(
		&p.ID, &p.Name, &p.Price, &p.CategoryID, &p.CreatedAt)

	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

func (r *ProductRepo) Create(p domain.Product) (domain.Product, error) {
	err := r.DB.QueryRow(
		context.Background(),
		"INSERT INTO products(name, price, category_id, created_at) VALUES($1, $2, $3, NOW()) RETURNING id, created_at",
		p.Name, p.Price, p.CategoryID).Scan(&p.ID, &p.CreatedAt)

	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

func (r *ProductRepo) Update(id int, p domain.Product) error {
	_, err := r.DB.Exec(
		context.Background(),
		"UPDATE products SET name = $1, price = $2, category_id = $3 WHERE id = $4",
		p.Name, p.Price, p.CategoryID, id,
	)
	return err
}

func (r *ProductRepo) Delete(id int) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM products WHERE id =$1", id)
	return err
}
