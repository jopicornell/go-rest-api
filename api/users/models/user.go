package models

import (
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	"strings"
)

const USER_ROLE = "user"
const ADMIN_ROLE = "admin"

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

// Returns true if user has the passed role
func (cwr *UserWithRoles) HasRole(roleToFind string) bool {
	roles := cwr.GetRoles()
	for _, role := range roles {
		if role == roleToFind {
			return true
		}
	}
	return false
}

// Returns true if user has all passed roles
func (cwr *UserWithRoles) HasSomeRole(rolesToFind []string) bool {
	roles := cwr.GetRoles()
	numRolesToMatch := len(rolesToFind)
	rolesMatched := 0
	for _, role := range rolesToFind {
		for _, roleToFind := range roles {
			if role == roleToFind {
				return true
			}
		}
	}
	return numRolesToMatch == rolesMatched
}

// Returns true if user has all passed roles
func (cwr *UserWithRoles) HasAllRoles(rolesToFind []string) bool {
	roles := cwr.GetRoles()
	numRolesToMatch := len(rolesToFind)
	rolesMatched := 0
	for _, role := range rolesToFind {
		for _, roleToFind := range roles {
			if role == roleToFind {
				numRolesToMatch++
			}
		}
	}
	return numRolesToMatch == rolesMatched
}

func (cwr *UserWithRoles) GetRolesString(separator string) string {
	return strings.Join(cwr.GetRoles(), separator)
}
