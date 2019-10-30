package tasks

import (
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/pkg/api/tasks/handlers"
)

func configureRoutes() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", handlers.GetTasksHandler)
	router.HandleFunc("/tasks/{id:[0-9]+}", handlers.GetOneTaskHandler)
	return
}
