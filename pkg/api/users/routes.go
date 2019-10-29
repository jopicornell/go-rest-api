package users

import (
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/pkg/api/users/handlers"
)

func configureRoutes() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlers.GetUsersHandler)
	return
}
