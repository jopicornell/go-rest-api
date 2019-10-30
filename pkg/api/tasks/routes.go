package tasks

import (
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/pkg/api"
	"github.com/jopicornell/go-rest-api/pkg/api/tasks/handlers"
)

func configureRoutes() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tasks", api.HandlerManager(handlers.GetTasksHandler))
	router.HandleFunc("/tasks/{id:[0-9]+}", api.HandlerManager(handlers.GetOneTaskHandler))
	return
}
