package infrastructure

import (
	"github.com/Lovodia/Product/internal/infrastructure/postgres"
	"github.com/Lovodia/Product/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryFactory struct {
	ProductRepo  usecase.ProductRepo
	CategoryRepo usecase.CategoryRepo
}

func NewRepositoryFactory(db *pgxpool.Pool) *RepositoryFactory {
	return &RepositoryFactory{
		ProductRepo:  postgres.NewProductRepo(db),
		CategoryRepo: postgres.NewCategoryRepo(db),
	}
}
