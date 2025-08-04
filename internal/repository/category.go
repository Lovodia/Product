package repository

import "github.com/Lovodia/Product/internal/domain"

type CategoryRepository interface {
	GetAll() ([]domain.Category, error)
	GetByID(id int) (domain.Category, error)
	Create(p domain.Category) (domain.Category, error)
	Update(category domain.Category) error
	Delete(id int) error
}
