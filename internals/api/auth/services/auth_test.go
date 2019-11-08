package services

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	dbMock, _ := mockDB(t)
	taskService := New(dbMock)
	if taskService == nil {
		t.Errorf("New service should not be null")
	}
	reflectServiceDb := reflect.ValueOf(taskService).Elem().FieldByName("db")
	if reflectServiceDb.IsNil() {
		t.Errorf("db field in service should not be nil")
	}
}

func TestAuthService_Login(t *testing.T) {
	t.Run("login should return unauthorized when passwords don't match", loginUnauthorized)
}

func loginUnauthorized(t *testing.T) {

}

func mockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}
