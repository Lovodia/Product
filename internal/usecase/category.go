package usecase

import (
	"github.com/Lovodia/Product/internal/domain"
)

type CategoryRepo interface {
	GetAll() ([]domain.Category, error)
	GetByID(id int) (domain.Category, error)
	Create(category domain.Category) (int, error)
	Update(category domain.Category) (bool, error)
	Delete(id int) (bool, error)
}

type CategoryUseCase struct {
	repo CategoryRepo
}

func NewCategoryUseCase(r CategoryRepo) *CategoryUseCase {
	return &CategoryUseCase{repo: r}
}

func (uc *CategoryUseCase) GetAll() ([]domain.Category, error) {
	return uc.repo.GetAll()
}

func (uc *CategoryUseCase) GetByID(id int) (domain.Category, error) {
	return uc.repo.GetByID(id)
}

func (uc *CategoryUseCase) Create(category domain.Category) (int, error) {
	return uc.repo.Create(category)
}

func (uc *CategoryUseCase) Update(category domain.Category) (bool, error) {
	return uc.repo.Update(category)
}

func (uc *CategoryUseCase) Delete(id int) (bool, error) {
	return uc.repo.Delete(id)
}
