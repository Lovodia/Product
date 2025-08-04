package usecase

import (
	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/repository"
)

type CategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewCategoryUseCase(r repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: r}
}

func (uc *CategoryUseCase) GetAll() ([]domain.Category, error) {
	return uc.repo.GetAll()
}

func (uc *CategoryUseCase) GetByID(id int) (domain.Category, error) {
	return uc.repo.GetByID(id)
}

func (uc *CategoryUseCase) Create(category domain.Category) (domain.Category, error) {
	created, err := uc.repo.Create(category)
	if err != nil {
		return domain.Category{}, err
	}
	return created, nil
}

func (uc *CategoryUseCase) Update(category domain.Category) error {
	return uc.repo.Update(category)
}

func (uc *CategoryUseCase) Delete(id int) error {
	return uc.repo.Delete(id)
}
