package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type MySQL struct {
	db  *sqlx.DB
	PSN string
}

func (m *MySQL) GetDB() *sqlx.DB {
	if m.db == nil {
		m.db = m.InitializeDB()
	}
	return m.db
}

func (m *MySQL) SetDB(dbInstance *sqlx.DB) {
	m.db = dbInstance
}

func (m *MySQL) InitializeDB() *sqlx.DB {
	var db *sqlx.DB
	err := Retry(func() (err error) {
		db, err = sqlx.Connect("mysql", m.PSN)
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
	return db
}
