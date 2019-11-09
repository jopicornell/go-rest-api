package services

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"golang.org/x/crypto/bcrypt"
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
	t.Run("should return error when user+password pair don't match", loginErrorNoMatch)
	t.Run("should return jwt when user+password pair match", loginReturnJWT)
}

func TestAuthService_Register(t *testing.T) {
	t.Run("should return error when weak password is provided", registerWeakPassword)
	t.Run("should return error when incorrect mail is provided", registerIncorrectMail)
	t.Run("should return error when email exists", registerMailExists)
	t.Run("should insert into database and return user when registered correctly", registerSuccessAndReturnUser)
}

func loginErrorNoMatch(t *testing.T) {
	mockedDb, mock := mockDB(t)
	authService := New(mockedDb, servertesting.Initialize(&config.Config{}))
	user := "user"
	password := "password"
	queryToRun := "SELECT id, name, email, active, deleted_at " +
		"FROM users WHERE email = \\?"
	mock.ExpectQuery(queryToRun).WillReturnRows(&sqlmock.Rows{})
	if got, err := authService.Login(user, password); err != nil {
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
	queryToRun := "SELECT id, name, email, active, deleted_at " +
		"FROM users WHERE email = \\?"
	rows := buildUserRows(true)
	_ = addUserRows(rows, []byte(password), 1)
	mock.ExpectQuery(queryToRun).WillReturnRows(rows)
	if got, err := authService.Login(user, password); err == nil {
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

func registerSuccessAndReturnUser(t *testing.T) {
	mockedDb, mock := mockDB(t)
	mockedServer := servertesting.Initialize(&config.Config{})
	authService := New(mockedDb, mockedServer)
	mock.ExpectBegin()
	insertQuery := "INSERT INTO users \\(name, email, password, active\\) VALUES " +
		"\\(\\?, \\?, \\?, \\?\\)"
	mock.ExpectExec(insertQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	rows := buildUserRows(false)
	users := addUserRows(rows, nil, 1)
	queryToRun := "SELECT id, name, email, active, deleted_at " +
		"FROM users WHERE id = LAST_INSERT_ID()"
	mock.ExpectQuery(queryToRun).WillReturnRows(rows)
	mock.ExpectCommit()
	if got, err := authService.Register(users[0]); err == nil {
		if !reflect.DeepEqual(*got, users[0]) {
			t.Errorf("expected return user %+v to be %+v", got, &users[0])
		}
	} else {
		t.Errorf("expected err %+v not to be nul", err)
	}
}

func registerMailExists(t *testing.T) {

}

func registerIncorrectMail(t *testing.T) {

}

func registerWeakPassword(t *testing.T) {

}

func mockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func buildUserRows(withPassword bool) *sqlmock.Rows {
	if withPassword {
		return sqlmock.NewRows([]string{"id", "name", "password", "email", "active", "deleted_at"})
	}
	return sqlmock.NewRows([]string{"id", "name", "email", "active", "deleted_at"})
}

func addUserRows(rows *sqlmock.Rows, password []byte, numRows uint) []models.User {
	var users []models.User
	var user *models.User
	for ; numRows > 0; numRows-- {
		user = createFakeUser(password)
		users = append(users, *user)
		if password == nil {
			rows.AddRow(user.ID, user.Name, user.Email, user.Active, user.DeletedAt)
		} else {
			rows.AddRow(user.ID, user.Name, user.Password, user.Email, user.Active, user.DeletedAt)
		}
	}
	return users
}

func createFakeUser(password []byte) (user *models.User) {
	user = &models.User{
		ID:     uint(faker.UnixTime()),
		Email:  faker.Email(),
		Name:   faker.Name(),
		Active: true,
		DeletedAt: null.Time{
			Time:  time.Now().Round(time.Nanosecond),
			Valid: true,
		},
	}
	if password != nil {
		user.Password = generatePasswordBcrypt(password)
	}
	return user
}

func generatePasswordBcrypt(password []byte) []byte {
	if password, err := bcrypt.GenerateFromPassword(password, 10); err == nil {
		return password
	} else {
		panic(err)
	}
}
