package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ManagesDatabases interface {
	GetRelationalDB() *sqlx.DB
}
