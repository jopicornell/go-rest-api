package users

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jopicornell/go-rest-api/pkg/util/config"
	"github.com/jopicornell/go-rest-api/pkg/util/database"
	"log"
	"net/http"
)

func Start() {
	config.Load()
	db := database.GetDB()
	if db != nil {
		defer db.Close()
	}
	router := configureRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
