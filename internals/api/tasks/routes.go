package tasks

import (
	"github.com/jopicornell/go-rest-api/internals/api/tasks/handlers"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

func ConfigureRoutes(s server.Server) {
	handler := handlers.New(s)
	s.AddApiRoute("/tasks", handler.GetTasksHandler, http.StatusOK).Methods("GET")
	s.AddApiRoute("/tasks", handler.CreateTaskHandler, http.StatusCreated).Methods("POST")
	s.AddApiRoute("/tasks", handler.DeleteTaskHandler, http.StatusCreated).Methods("DELETE")
	s.AddApiRoute("/tasks/{id:[0-9]+}", handler.GetOneTaskHandler, http.StatusOK).Methods("GET")
	s.AddApiRoute("/tasks/{id:[0-9]+}", handler.UpdateTaskHandler, http.StatusOK).Methods("PUT")
	s.AddStatics("/", "static")
}
