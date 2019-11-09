package models

import (
	"gopkg.in/guregu/null.v3"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"  db:"name"`
	Password  []byte    `db:"password"`
	Email     string    `json:"email"  db:"email"`
	Active    bool      `json:"active"  db:"active"`
	CreatedAt null.Time `json:"created_at" db:"created_at"`
	UpdatedAt null.Time `json:"updated_at" db:"updated_at"`
	DeletedAt null.Time `json:"deleted_at" db:"deleted_at"`
}
