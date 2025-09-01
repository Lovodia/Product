package postgres_test

import (
	"testing"

	"github.com/Lovodia/Product/internal/domain"
	"github.com/Lovodia/Product/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepo_GetAll(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)

	rows := pgxmock.NewRows([]string{"id", "name"}).
		AddRow(1, "A").AddRow(2, "B")
	mock.ExpectQuery("SELECT id, name FROM categories").
		WillReturnRows(rows)

	got, err := repo.GetAll(t.Context())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, domain.Category{ID: 1, Name: "A"}, got[0])
	require.Equal(t, domain.Category{ID: 2, Name: "B"}, got[1])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_GetByID_Success(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)
	expected := domain.Category{ID: 1, Name: "Electronics"}
	mock.ExpectQuery("SELECT id, name FROM categories WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name"}).
			AddRow(expected.ID, expected.Name))

	got, err := repo.GetByID(t.Context(), 1)
	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)
	mock.ExpectQuery("SELECT id, name FROM categories WHERE id=\\$1").
		WithArgs(10).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByID(t.Context(), 10)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Create(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)
	category := domain.Category{Name: "Boocs"}
	mock.ExpectQuery("INSERT INTO categories\\(name\\) VALUES\\(\\$1\\) RETURNING id").
		WithArgs(category.Name).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(10))

	id, err := repo.Create(t.Context(), category)
	require.NoError(t, err)
	require.Equal(t, 10, id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Update_Success(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)
	category := domain.Category{ID: 5, Name: "Updated"}
	mock.ExpectExec("UPDATE categories SET name=\\$1 WHERE id=\\$2").
		WithArgs(category.Name, category.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	ok, err := repo.Update(t.Context(), category)
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Update_NoFound(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)
	category := domain.Category{ID: 99, Name: "Nope"}
	mock.ExpectExec("UPDATE categories SET name=\\$1 WHERE id=\\$2").
		WithArgs(category.Name, category.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	ok, err := repo.Update(t.Context(), category)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Delete_Success(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)
	mock.ExpectExec("DELETE FROM categories WHERE id=\\$1").
		WithArgs(3).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	ok, err := repo.Delete(t.Context(), 3)
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Delete_NotFound(t *testing.T) {
	t.Parallel()

	mock, _ := pgxmock.NewPool()

	defer mock.Close()

	repo := postgres.NewCategoryRepo(mock)
	mock.ExpectExec("DELETE FROM categories WHERE id=\\$1").
		WithArgs(404).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	ok, err := repo.Delete(t.Context(), 404)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}
