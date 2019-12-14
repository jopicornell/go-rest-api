package models

import (
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	"strings"
)

const USER_ROLE = "user"

type User struct {
	model.User
}

type UserWithRoles struct {
	model.User
	Roles []model.UserHasRoles
}

func (cwr *UserWithRoles) GetRoles() []string {
	roles := make([]string, 0)
	for _, role := range cwr.Roles {
		roles = append(roles, role.Role)
	}
	return roles
}

func (cwr *UserWithRoles) GetRolesString(separator string) string {
	return strings.Join(cwr.GetRoles(), separator)
}
