package usecase

import (
	"context"

	"github.com/Lovodia/Product/internal/domain"
)

type ProductRepo interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id int) (domain.Product, error)
	Create(ctx context.Context, p domain.Product) (int, error)
	Update(ctx context.Context, id int, p domain.Product) (bool, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type ProductUseCase struct {
	repo ProductRepo
}

func NewProductUseCase(r ProductRepo) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

func (uc *ProductUseCase) GetAll(ctx context.Context) ([]domain.Product, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *ProductUseCase) GetByID(ctx context.Context, id int) (domain.Product, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ProductUseCase) Create(ctx context.Context, p domain.Product) (int, error) {
	return uc.repo.Create(ctx, p)
}

func (uc *ProductUseCase) Update(ctx context.Context, id int, p domain.Product) (bool, error) {
	return uc.repo.Update(ctx, id, p)
}

func (uc *ProductUseCase) Delete(ctx context.Context, id int) (bool, error) {
	return uc.repo.Delete(ctx, id)
}
