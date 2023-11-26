package db

import (
	"database/sql"
	"errors"
	"example_mock/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShouldGetNames(t *testing.T) {
	testTable := []struct {
		dbRows       *sqlmock.Rows
		expectedRows []string
		expectedErr  error
	}{
		{
			dbRows: sqlmock.NewRows([]string{"name"}).
				AddRow("Khamzat").
				AddRow("Khabib").
				AddRow("Merab"),
			expectedRows: []string{"Khamzat", "Khabib", "Merab"},
			expectedErr:  nil,
		},
		{
			dbRows: sqlmock.NewRows([]string{"name"}).
				AddRow("Yan").
				AddRow("Borya").
				RowError(1, errors.New("row error")),
			expectedRows: nil,
			expectedErr:  errors.New("row error"),
		},
		{
			dbRows: sqlmock.NewRows([]string{"name"}).
				AddRow("Khamzat").
				AddRow("Khabib").
				AddRow("Merab").
				CloseError(errors.New("close error")),
			expectedRows: nil,
			expectedErr:  errors.New("close error"),
		},
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	for _, row := range testTable {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(row.dbRows)
		names, err := dbService.GetNames()

		if row.expectedErr != nil {
			require.Error(t, err)
			require.Nil(t, names)
			continue
		}

		require.NoError(t, err)
		require.Equal(t, row.expectedRows, names)
	}
}

func TestShouldGetUniqueNames(t *testing.T) {
	testTable := []struct {
		dbRows       *sqlmock.Rows
		expectedRows []string
		expectedErr  error
	}{
		{
			dbRows: sqlmock.NewRows([]string{"name"}).
				AddRow("Khamzat").
				AddRow("Khabib"),
			expectedRows: []string{"Khamzat", "Khabib"},
			expectedErr:  nil,
		},
		{
			dbRows: sqlmock.NewRows([]string{"name"}).
				AddRow("Yan").
				AddRow("Borya").
				RowError(1, errors.New("row error")),
			expectedRows: nil,
			expectedErr:  errors.New("row error"),
		},
		{
			dbRows: sqlmock.NewRows([]string{"name"}).
				CloseError(sql.ErrNoRows),
			expectedRows: nil,
			expectedErr:  sql.ErrNoRows,
		},
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	for _, row := range testTable {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(row.dbRows)
		names, err := dbService.SelectUniqueValues("name", "users")

		if row.expectedErr != nil {
			require.Error(t, err)
			require.Nil(t, names)
			continue
		}

		require.NoError(t, err)
		require.Equal(t, row.expectedRows, names)
	}
}
