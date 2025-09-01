package infrastructure

import (
	"github.com/Lovodia/Product/internal/infrastructure/postgres"
	"github.com/Lovodia/Product/internal/usecase"
)

type RepositoryFactory struct {
	ProductRepo  usecase.ProductRepo
	CategoryRepo usecase.CategoryRepo
}

func NewRepositoryFactory(db postgres.PgxIface) *RepositoryFactory {
	return &RepositoryFactory{
		ProductRepo:  postgres.NewProductRepo(db),
		CategoryRepo: postgres.NewCategoryRepo(db),
	}
}
