package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/models"
)

type AuthService interface {
	Login() ([]models.Token, error)
}

type authService struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) AuthService {
	return &authService{
		db: db,
	}
}

func (s *authService) Login() (tasks []models.Token, err error) {
	return nil, nil
}
