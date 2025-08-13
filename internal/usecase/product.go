package usecase

import (
	"github.com/Lovodia/Product/internal/domain"
)

type ProductRepo interface {
	GetAll() ([]domain.Product, error)
	GetByID(id int) (domain.Product, error)
	Create(p domain.Product) (int, error)
	Update(id int, p domain.Product) (bool, error)
	Delete(id int) (bool, error)
}

type ProductUseCase struct {
	repo ProductRepo
}

func NewProductUseCase(r ProductRepo) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

func (uc *ProductUseCase) GetAll() ([]domain.Product, error) {
	return uc.repo.GetAll()
}

func (uc *ProductUseCase) GetByID(id int) (domain.Product, error) {
	return uc.repo.GetByID(id)
}

func (uc *ProductUseCase) Create(p domain.Product) (int, error) {
	return uc.repo.Create(p)
}

func (uc *ProductUseCase) Update(id int, p domain.Product) (bool, error) {
	return uc.repo.Update(id, p)
}

func (uc *ProductUseCase) Delete(id int) (bool, error) {
	return uc.repo.Delete(id)
}
