package services

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"gopkg.in/guregu/null.v3"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	dbMock, _ := mockDB(t)
	taskService := New(dbMock, servertesting.Initialize(&config.Config{}))
	if taskService == nil {
		t.Errorf("New service should not be null")
	}
	reflectServiceDb := reflect.ValueOf(taskService).Elem().FieldByName("db")
	if reflectServiceDb.IsNil() {
		t.Errorf("db field in service should not be nil")
	}
}

func TestAuthService_Login(t *testing.T) {
	t.Run("login should return error when user+password pair don't match", loginErrorNoMatch)
	t.Run("login should return jwt when user+password pair match", loginReturnJWT)
}

func loginErrorNoMatch(t *testing.T) {
	mockedDb, mock := mockDB(t)
	authService := New(mockedDb, servertesting.Initialize(&config.Config{}))
	user := "user"
	password := "password"
	queryToRun := "SELECT id, name, email, user_name, active, deleted_at " +
		"FROM users WHERE sha2\\(concat\\(first_name, last_name\\), 256\\) = \\?"
	mock.ExpectQuery(queryToRun).WillReturnRows(&sqlmock.Rows{})
	hasher := sha256.New()
	tokenizedRequest := base64.URLEncoding.EncodeToString(hasher.Sum([]byte(user + password)))
	if got, err := authService.Login(string(tokenizedRequest)); err != nil {
		if got != nil {
			t.Errorf("expected return %+v to be nil", got)
		}
		if err != errors.AuthUserNotMatched {
			t.Errorf("expected %v to be %v", err, errors.AuthUserNotMatched)
		}
	} else {
		t.Errorf("expected to throw error to be nil")
	}
}

func loginReturnJWT(t *testing.T) {
	mockedDb, mock := mockDB(t)
	mockedServer := servertesting.Initialize(&config.Config{})
	authService := New(mockedDb, mockedServer)
	user := "user"
	password := "password"
	queryToRun := "SELECT id, name, email, user_name, active, deleted_at " +
		"FROM users WHERE sha2\\(concat\\(first_name, last_name\\), 256\\) = \\?"
	rows := buildUserRows()
	_ = addUserRows(rows, 1)
	mock.ExpectQuery(queryToRun).WillReturnRows(rows)
	hasher := sha256.New()
	tokenizedRequest := base64.URLEncoding.EncodeToString(hasher.Sum([]byte(user + password)))
	if got, err := authService.Login(string(tokenizedRequest)); err == nil {
		if got == nil {
			t.Errorf("expected return not to be nil")
		}
		if err := validateJWT(got.Token, mockedServer.GetServerConfig().JWTSecret); err != nil {
			t.Errorf("token could not be verified correctly:  %w", err)
		}
	} else {
		t.Errorf("expected to not throw error, throwed %w", err)
	}
}

func validateJWT(token string, secret string) error {
	var hs = jwt.NewHS512([]byte(secret))
	var payload *models.JwtUserPayload
	_, err := jwt.Verify([]byte(token), hs, &payload)
	return err
}

func mockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func buildUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "email", "user_name", "active", "deleted_at"})
}

func addUserRows(rows *sqlmock.Rows, numRows uint) []models.User {
	var users []models.User
	var user *models.User
	for ; numRows > 0; numRows-- {
		user = createFakeUser()

		users = append(users, *user)
		rows.AddRow(user.ID, user.Name, user.Email, user.UserName, user.Active, user.DeletedAt)
	}
	return users
}

func createFakeUser() *models.User {
	return &models.User{
		ID:       uint(faker.UnixTime()),
		Email:    faker.Email(),
		Name:     faker.Name(),
		UserName: faker.Username(),
		Active:   true,
		DeletedAt: null.Time{
			Time:  time.Now().Round(time.Nanosecond),
			Valid: true,
		},
	}
}
