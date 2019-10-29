package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/util/database"
	"github.com/jopicornell/go-rest-api/pkg/util/models"
)

type Service struct {
	DB *sqlx.DB
}

func (s *Service) GetUsers() (users []models.User, err error) {
	err = database.GetDB().Select(&users, "SELECT * from users")
	return
}
