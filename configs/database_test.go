//go:build unit

package configs

import (
	"database/sql"
	"errors"
	"testing"
)

type mockDB struct {
	query        string
	lastInsertId int64
	rowsAffected int64
}

// Exec is a function that execute query
func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.query = query
	return m, nil
}

// LastInsertId is a function that return last insert id
func (m *mockDB) LastInsertId() (int64, error) {
	return m.lastInsertId, nil
}

// RowsAffected is a function that return rows affected
func (m *mockDB) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func TestAutoMigrate(t *testing.T) {
	mock := &mockDB{}

	AutoMigrate(mock)

	if mock.query != `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	` {
		t.Error("should have been call db.Exec with but it's not")
	}
}

func TestOpenDB(t *testing.T) {
	mockError := errors.New("Mock Error")
	subtests := []struct {
		name          string
		connectionUrl string
		sqlOpener     func(string, string) (*sql.DB, error)
		expectedErr   error
	}{
		{
			name:          "test wrong path",
			connectionUrl: "mock path",
			sqlOpener: func(s string, s2 string) (db *sql.DB, err error) {
				if s != "postgres" {
					return nil, errors.New("wrong database type")
				}
				if s2 != "mock path" {
					return nil, errors.New("wrong database connection url")
				}
				return nil, nil
			},
		},
		{
			name:          "test db error",
			connectionUrl: "mock path",
			sqlOpener: func(s string, s2 string) (db *sql.DB, err error) {
				return nil, mockError
			},
			expectedErr: mockError,
		},
	}
	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			_, err := OpenDB(subtest.sqlOpener, subtest.connectionUrl)
			if !errors.Is(err, subtest.expectedErr) {
				t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
			}
		})
	}
}
