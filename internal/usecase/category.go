package usecase

import (
	"context"

	"github.com/Lovodia/Product/internal/domain"
)

type CategoryRepo interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id int) (domain.Category, error)
	Create(ctx context.Context, category domain.Category) (int, error)
	Update(ctx context.Context, category domain.Category) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type CategoryUseCase struct {
	repo CategoryRepo
}

func NewCategoryUseCase(r CategoryRepo) *CategoryUseCase {
	return &CategoryUseCase{repo: r}
}

func (uc *CategoryUseCase) GetAll(ctx context.Context) ([]domain.Category, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *CategoryUseCase) GetByID(ctx context.Context, id int) (domain.Category, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *CategoryUseCase) Create(ctx context.Context, category domain.Category) (int, error) {
	return uc.repo.Create(ctx, category)
}

func (uc *CategoryUseCase) Update(ctx context.Context, category domain.Category) (bool, error) {
	return uc.repo.Update(ctx, category)
}

func (uc *CategoryUseCase) Delete(ctx context.Context, id int) (bool, error) {
	return uc.repo.Delete(ctx, id)
}
