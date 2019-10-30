package tasks

import (
	"github.com/jopicornell/go-rest-api/pkg/database"
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
