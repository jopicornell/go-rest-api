package services

import (
	"database/sql"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"strconv"
)

type AuthService interface {
	Login(tokenizedRequest string) (*models.Token, error)
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

func (s *authService) Login(tokenizedRequest string) (token *models.Token, err error) {
	query := "SELECT id, name, email, user_name, active, deleted_at FROM users " +
		"WHERE sha2(concat(first_name, last_name), 256) = ?"
	var user models.User
	if err := s.db.Get(&user, query, tokenizedRequest); err == nil {
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
