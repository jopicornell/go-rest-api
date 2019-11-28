package database

import (
	"github.com/jmoiron/sqlx"
)

type ManagesDatabases interface {
	GetRelationalDB() *sqlx.DB
}
