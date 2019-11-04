package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
)

type MySQL struct {
	db  *sqlx.DB
	PSN string
}

func (m *MySQL) GetDB() *sqlx.DB {
	if m.db == nil {
		m.InitializeDB()
	}
	return m.db
}

func (m *MySQL) SetDB(dbInstance *sqlx.DB) {
	m.db = dbInstance
}

func (m *MySQL) InitializeDB() {
	var err error
	m.db, err = sqlx.Connect("mysql", m.PSN)
	if err != nil {
		wrapError := errors.Wrap(err, "some problem with initializing m")
		log.Printf(wrapError.Error())
	}
}
