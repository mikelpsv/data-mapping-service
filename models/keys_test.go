package models

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mikelpsv/data-mapping-service/app"
	"log"
	"testing"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestKey_FindById_Ok(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	app.Db = db

	rows := sqlmock.NewRows([]string{"_id", "name"})
	rows.AddRow(1, "Ключ 1")

	mock.ExpectQuery("^SELECT (.+) FROM keys*").
		WillReturnRows(rows)

	key := new(Key)
	key.Id = 10

	key, err := key.FindById(1)
	if err != nil {
		t.Error(err)
	}
	if key.Id != 1 {
		t.Error("Result ID fail")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestKey_FindById_InvalidKey(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	app.Db = db

	rows := sqlmock.NewRows([]string{"_id", "name"})
	mock.ExpectQuery("^SELECT (.+) FROM keys*").
		WillReturnRows(rows)

	key := new(Key)
	key.Id = 10
	key.Name = "Test"

	_, err := key.FindById(3)
	if err == nil {
		t.Error("Error is nil")
	}

	if key.Id != 0 || key.Name != "" {
		t.Error("Prev object data fail")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
