package database

import (
	"github.com/jmoiron/sqlx"
)

type ManagesDatabases interface {
	GetRelationalDatabase() *sqlx.DB
	GetCache() Cache
}
