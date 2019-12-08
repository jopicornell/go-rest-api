package database

import (
	_ "github.com/jackc/pgx/stdlib"
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
	return m.db
}

func (m *Postgres) SetDB(dbInstance *sqlx.DB) {
	m.db = dbInstance
}

func (m *Postgres) InitializeDB() {
	var db *sqlx.DB
	err := Retry(func() (err error) {
		db, err = sqlx.Open("pgx", m.PSN)
		if err != nil {
			logrus.Errorf("Error connecting to client %s", err)
		}
		return err
	}, time.Second*15, time.Minute*5)
	if err != nil {
		wrapError := errors.Wrap(err, "some problem with initializing relational client")
		log.Fatal(wrapError.Error())
	}
	log.Println("New connection to postgres")
	db.SetMaxOpenConns(10)
	m.db = db
}

func InitializeDB(PSN string) *sqlx.DB {
	var db *sqlx.DB
	err := Retry(func() (err error) {
		db, err = sqlx.Open("pgx", PSN)
		if err != nil {
			logrus.Errorf("Error connecting to client %s", err)
		}
		return err
	}, time.Second*15, time.Minute*5)
	if err != nil {
		wrapError := errors.Wrap(err, "some problem with initializing relational client")
		log.Fatal(wrapError.Error())
	}
	log.Println("New connection to postgres")
	db.SetMaxOpenConns(10)
	return db
}
