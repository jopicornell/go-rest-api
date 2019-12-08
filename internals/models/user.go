package models

import (
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	"strings"
)

const USER_ROLE = "user"

type User struct {
	model.Customer
}

type CustomerWithRoles struct {
	model.Customer
	Roles []model.CustomerHasRoles
}

func (cwr *CustomerWithRoles) GetRoles() []string {
	roles := make([]string, 0)
	for _, role := range cwr.Roles {
		roles = append(roles, role.Role)
	}
	return roles
}

func (cwr *CustomerWithRoles) GetRolesString(separator string) string {
	return strings.Join(cwr.GetRoles(), separator)
}
