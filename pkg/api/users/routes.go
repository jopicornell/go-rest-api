package users

import "github.com/gorilla/mux"

func configureRoutes() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", getUsersHandler)
	return
}
