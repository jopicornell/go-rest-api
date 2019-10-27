package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/util/config"
	"log"
)

var db *sqlx.DB = nil

func GetDB() *sqlx.DB {
	if db == nil {
		db = startDB()
	}
	return db
}

func startDB() *sqlx.DB {
	var err error
	db, err = sqlx.Connect("mysql", config.GetDBConfig().PSN)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
