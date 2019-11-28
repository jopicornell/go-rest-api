package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type Postgres struct {
	db  *sqlx.DB
	PSN string
}

func (m *Postgres) GetDB() *sqlx.DB {
	if m.db == nil {
		m.db = m.InitializeDB()
	}
	return m.db
}

func (m *Postgres) SetDB(dbInstance *sqlx.DB) {
	m.db = dbInstance
}

func (m *Postgres) InitializeDB() *sqlx.DB {
	var db *sqlx.DB
	err := Retry(func() (err error) {
		db, err = sqlx.Open("postgres", m.PSN)
		if err != nil {
			logrus.Errorf("Error connecting to db %s", err)
		}
		return err
	}, time.Second*15, time.Minute*5)
	if err != nil {
		wrapError := errors.Wrap(err, "some problem with initializing relational db")
		log.Fatal(wrapError.Error())
	}
	log.Println("New connection to db")
	db.SetMaxOpenConns(10)
	return db
}
