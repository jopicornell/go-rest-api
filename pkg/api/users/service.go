package users

import (
	"github.com/jopicornell/go-rest-api/pkg/util/database"
	"github.com/jopicornell/go-rest-api/pkg/util/models"
)

type Service interface {
	getUsers() (users []models.User, err error)
}

func getUsers() (users []models.User, err error) {
	err = database.GetDB().Select(&users, "SELECT * from users")
	return
}
