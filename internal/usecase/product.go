package usecase

import (
	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/repository"
)

type ProductUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(r repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

func (uc *ProductUseCase) GetAll() ([]domain.Product, error) {
	return uc.repo.GetAll()
}

func (uc *ProductUseCase) GetByID(id int) (domain.Product, error) {
	return uc.repo.GetByID(id)
}

func (uc *ProductUseCase) Create(p domain.Product) (domain.Product, error) {
	return uc.repo.Create(p)
}

func (uc *ProductUseCase) Update(id int, p domain.Product) error {
	return uc.repo.Update(id, p)
}

func (uc *ProductUseCase) Delete(id int) error {
	return uc.repo.Delete(id)
}
