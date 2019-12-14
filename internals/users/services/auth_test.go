package services

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/config"
	password2 "github.com/jopicornell/go-rest-api/pkg/password"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"reflect"
	"testing"
)

func TestNewAuthService(t *testing.T) {
	dbMock, _ := mockAuthDB(t)
	taskService := NewAuthService(dbMock, servertesting.Initialize(&config.Config{}))
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
	mockedDb, mock := mockAuthDB(t)
	authService := NewAuthService(mockedDb, servertesting.Initialize(&config.Config{}))
	user := "user"
	password := "password"
	mock.ExpectQuery(".*").WillReturnRows(&sqlmock.Rows{})
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
	mockedDb, mock := mockAuthDB(t)
	mockedServer := servertesting.Initialize(&config.Config{})
	authService := NewAuthService(mockedDb, mockedServer)
	user := "user"
	password := "password"
	rows := buildUserRows()
	_ = addUserRows(rows, password, 1)
	mock.ExpectQuery(".*").WillReturnRows(rows)
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
	mockedDb, mock := mockAuthDB(t)
	mockedServer := servertesting.Initialize(&config.Config{})
	authService := NewAuthService(mockedDb, mockedServer)
	mock.ExpectBegin()
	insertQuery := "INSERT INTO users \\(name, email, password, active\\) VALUES " +
		"\\(\\?, \\?, \\?, \\?\\)"
	mock.ExpectExec(insertQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	rows := buildUserRows()
	users := addUserRows(rows, "password", 1)
	queryToRun := "SELECT id, name, email, active, deleted_at " +
		"FROM users WHERE id = LAST_INSERT_ID()"
	mock.ExpectQuery(queryToRun).WillReturnRows(rows)
	mock.ExpectCommit()
	if got, err := authService.Register(&users[0]); err == nil {
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

func mockAuthDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func buildUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"customer.user_id", "customer.username", "customer.full_name", "customer.password"})
}

func addUserRows(rows *sqlmock.Rows, password string, numRows uint) []model.User {
	var users []model.User
	var user *model.User
	for ; numRows > 0; numRows-- {
		user = createFakeUser(password)
		users = append(users, *user)
		rows.AddRow(user.UserID, user.Username, user.FullName, user.Password)
	}
	return users
}

func createFakeUser(password string) (user *model.User) {
	user = &model.User{
		UserID:   int32(faker.UnixTime()),
		FullName: faker.Name(),
		Password: generatePasswordArgon(password),
	}
	return user
}

func generatePasswordArgon(passwd string) string {
	if password, err := password2.ArgonHashFromPassword(passwd, &password2.ArgonPasswordParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}); err == nil {
		return password
	} else {
		panic(err)
	}
}
