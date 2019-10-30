package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrationDriver "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"log"
	"strconv"
)

func main() {
	var err error
	var driver migrationDriver.Driver
	driver, err = mysql.WithInstance(database.GetDB().DB, &mysql.Config{})
	if err != nil {
		log.Fatal("Unable to connect to db")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql", driver)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to run migrations: %w", err))
	}
	flag.Parse()
	command := flag.Arg(0)
	switch command {
	case "up":
		if err = m.Up(); err != nil {
			handleMigrationErr(err)
		}
	case "down":
		if err = m.Down(); err != nil {
			handleMigrationErr(err)
		}

	case "goto":
		if flag.Arg(1) == "" {
			log.Fatal("error: please specify version argument V")
		}
		version, err := strconv.ParseUint(flag.Arg(1), 10, 64)
		if err != nil {
			log.Fatal("error: can't read version")
		}
		if err = m.Migrate(uint(version)); err != nil {
			handleMigrationErr(err)
		}
	default:
		if err = m.Up(); err != nil {
			handleMigrationErr(err)
		}
	}

	fmt.Println("Everything executed correctly")
}

func handleMigrationErr(err error) {
	if err == migrate.ErrNoChange {
		log.Println(err)
	} else {
		log.Fatal(err)
	}
}
