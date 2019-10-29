package migrate

import (
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/golang-migrate/migrate/v4"
	migrationDriver "github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/jopicornell/go-rest-api/pkg/util/database"
	"log"
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
		"postgres", driver)
	if err != nil {
		log.Fatalf("Unable to run migrations %")
	}
	m.Up()
}
