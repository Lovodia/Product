package postgres_test

import (
	"testing"
	"time"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestProductRepo_GetAll(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)

	rows := pgxmock.NewRows([]string{"id", "name", "price", "category_id", "created_at"}).
		AddRow(1, "P1", 9.9, 2, time.Now()).
		AddRow(2, "P2", 19.9, 3, time.Now())
	mock.ExpectQuery("SELECT id, name, price, category_id, created_at FROM products").
		WillReturnRows(rows)

	got, err := repo.GetAll(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, 1, got[0].ID)
	require.Equal(t, 2, got[1].ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_GetByID_Success(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)
	now := time.Now()
	expected := domain.Product{ID: 7, Name: "X", Price: 5.5, CategoryID: 1, CreatedAt: now}

	mock.ExpectQuery("SELECT id, name, price, category_id, created_at FROM products WHERE id = \\$1").
		WithArgs(7).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "price", "category_id", "created_at"}).
			AddRow(expected.ID, expected.Name, expected.Price, expected.CategoryID, now))

	prod, err := repo.GetByID(t.Context(), 7)
	require.NoError(t, err)
	require.Equal(t, expected, prod)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)
	mock.ExpectQuery("SELECT id, name, price, category_id, created_at FROM products WHERE id = \\$1").
		WithArgs(123).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(t.Context(), 123)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_Create(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)
	prod := domain.Product{Name: "New", Price: 3.3, CategoryID: 2}
	mock.ExpectQuery(
		"INSERT INTO products\\(name, price, category_id, created_at\\) VALUES\\(\\$1, \\$2, \\$3, NOW\\(\\)\\) RETURNING id, created_at").
		WithArgs(prod.Name, prod.Price, prod.CategoryID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at"}).
			AddRow(11, time.Now()))

	id, err := repo.Create(t.Context(), prod)
	require.NoError(t, err)
	require.Equal(t, 11, id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_Update_Success(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)
	prod := domain.Product{Name: "Upd", Price: 4.4, CategoryID: 3}
	mock.ExpectExec("UPDATE products SET name = \\$1, price = \\$2, category_id = \\$3 WHERE id = \\$4").
		WithArgs(prod.Name, prod.Price, prod.CategoryID, 5).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	ok, err := repo.Update(t.Context(), 5, prod)
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_Update_NotFound(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)
	prod := domain.Product{Name: "Missing", Price: 0, CategoryID: 0}
	mock.ExpectExec("UPDATE products SET name = \\$1, price = \\$2, category_id = \\$3 WHERE id = \\$4").
		WithArgs(prod.Name, prod.Price, prod.CategoryID, 404).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	ok, err := repo.Update(t.Context(), 404, prod)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_Delete_Success(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)
	mock.ExpectExec("DELETE FROM products WHERE id =\\$1").
		WithArgs(8).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	ok, err := repo.Delete(t.Context(), 8)
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestProductRepo_Delete_NotFound(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewProductRepo(mock)
	mock.ExpectExec("DELETE FROM products WHERE id =\\$1").
		WithArgs(999).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	ok, err := repo.Delete(t.Context(), 999)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}
