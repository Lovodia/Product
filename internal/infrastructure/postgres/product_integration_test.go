package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Lovodia/Product/internal/config"
	"github.com/Lovodia/Product/internal/db"
	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/infrastructure/postgres"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupProductRepo(t *testing.T) (*postgres.ProductRepo, *postgres.CategoryRepo, func()) {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := &config.Config{
		DB: config.DBConfig{
			Host:     host,
			Port:     port.Int(),
			User:     "user",
			Password: "pass",
			Name:     "testdb",
			SSLMode:  "disable",
		},
		Migration: config.MigrationConfig{
			Path: "../../../migrations",
		},
	}

	database, err := db.New(cfg)
	require.NoError(t, err)

	err = db.RunMigrations(database.Pool(), cfg.Migration.Path)
	require.NoError(t, err)

	productRepo := postgres.NewProductRepo(database.Pool())
	categoryRepo := postgres.NewCategoryRepo(database.Pool())

	return productRepo, categoryRepo, func() {
		database.Close()
		_ = container.Terminate(ctx)
	}
}

func TestProductRepo_CRUD(t *testing.T) {
	t.Parallel()

	productRepo, categoryRepo, cleanup := setupProductRepo(t)
	defer cleanup()

	ctx := context.Background()

	catID, err := categoryRepo.Create(ctx, domain.Category{Name: "Tech"})
	require.NoError(t, err)

	prod := domain.Product{
		Name:       "Phone",
		Price:      999.99,
		CategoryID: catID,
	}

	id, err := productRepo.Create(ctx, prod)
	require.NoError(t, err)
	require.True(t, id > 0)

	p, err := productRepo.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "Phone", p.Name)

	p.Name = "Smartphone"
	p.Price = 899.99
	ok, err := productRepo.Update(ctx, id, p)
	require.NoError(t, err)
	require.True(t, ok)

	all, err := productRepo.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)

	ok, err = productRepo.Delete(ctx, id)
	require.NoError(t, err)
	require.True(t, ok)

	_, err = productRepo.GetByID(ctx, id)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
