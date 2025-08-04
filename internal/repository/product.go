package repository

import "github.com/Lovodia/Product/internal/domain"

type ProductRepository interface {
	GetAll() ([]domain.Product, error)
	GetByID(id int) (domain.Product, error)
	Create(p domain.Product) (domain.Product, error)
	Update(id int, p domain.Product) error
	Delete(id int) error
}
