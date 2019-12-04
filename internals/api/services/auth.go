package services

import (
	"crypto/sha512"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Login(email string, password string) (*models.Token, error)
	Register(user *models.User) (*models.User, error)
}

type authService struct {
	db     *sqlx.DB
	server server.Server
}

func NewAuthService(db *sqlx.DB, server server.Server) AuthService {
	return &authService{
		db:     db,
		server: server,
	}
}

func (s *authService) Login(email string, password string) (token *models.Token, err error) {
	query := "SELECT id, name, password, email, active, deleted_at FROM auth.users " +
		"WHERE email = $1"
	var user models.User
	if err := s.db.Get(&user, query, email); err == nil {
		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
			return nil, errors.AuthUserNotMatched
		}
		token := sha512.New().Sum([]byte(fmt.Sprintf("%d;%s;%d", user.ID, user.Password, time.Now().Nanosecond())))
		return &models.Token{Token: string(token)}, nil
	} else {
		if err == sql.ErrNoRows {
			return nil, errors.AuthUserNotMatched
		}
		return nil, err
	}
}

func (s *authService) Register(user *models.User) (_ *models.User, err error) {
	tx := s.db.MustBegin()
	insertStatement := "INSERT INTO auth.users (name, email, password, active) VALUES ($1, $2, $3, $4) RETURNING id;"
	if result := tx.QueryRowx(insertStatement, user.Name, user.Email, user.Password, true); result.Err() == nil {
		if err := result.Scan(&user.ID); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		if err = tx.Commit(); err == nil {
			user.Password = nil
			return user, nil
		}
	}
	return nil, err
}
