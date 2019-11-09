package services

import (
	"database/sql"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthService interface {
	Login(email string, password string) (*models.Token, error)
	Register(user models.User) (*models.User, error)
}

type authService struct {
	db     *sqlx.DB
	server server.Server
}

func New(db *sqlx.DB, server server.Server) AuthService {
	return &authService{
		db:     db,
		server: server,
	}
}

func (s *authService) Login(email string, password string) (token *models.Token, err error) {
	query := "SELECT id, name, email, active, deleted_at FROM users " +
		"WHERE email = ?"
	var user models.User
	if err := s.db.Get(&user, query, email); err == nil {
		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
			return nil, errors.AuthUserNotMatched
		}
		var hs = jwt.NewHS512([]byte(s.server.GetServerConfig().JWTSecret))
		token, err := jwt.Sign(configurePayload(user.ID), hs)
		if err != nil {
			return nil, err
		}
		return &models.Token{Token: string(token)}, nil
	} else {
		if err == sql.ErrNoRows {
			return nil, errors.AuthUserNotMatched
		}
		return nil, err
	}
}

func (s *authService) Register(user models.User) (_ *models.User, err error) {
	tx := s.db.MustBegin()
	insertStatement := "INSERT INTO users (name, email, password, active) VALUES (?, ?, ?, ?)"
	if _, err = tx.Exec(insertStatement, user.Name, user.Email, user.Password, true); err == nil {
		queryToRun := "SELECT id, name, email, active, deleted_at " +
			"FROM users WHERE id = LAST_INSERT_ID()"
		var user models.User
		if err = tx.Get(&user, queryToRun); err == nil {
			if err = tx.Commit(); err == nil {
				return &user, nil
			}
		}
	}
	return nil, err
}

func configurePayload(userId uint) models.JwtUserPayload {
	return models.JwtUserPayload{
		Payload: jwt.Payload{
			Issuer:         "palmaactiva",
			Subject:        strconv.Itoa(int(userId)),
			Audience:       nil,
			ExpirationTime: nil,
			NotBefore:      nil,
			IssuedAt:       nil,
			JWTID:          "",
		},
		ID: int(userId),
	}
}
