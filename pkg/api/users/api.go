package users

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jopicornell/go-rest-api/pkg/util/database"
	"log"
	"net/http"
)

func Start() {
	db := database.GetDB()
	if db != nil {
		defer db.Close()
	}
	router := configureRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
