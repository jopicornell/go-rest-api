package tasks

import (
	"github.com/jopicornell/go-rest-api/internals/api/tasks/handlers"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s *server.Server) {
	handler := handlers.New(s)
	s.AddApiRoute("/tasks", handler.GetTasksHandler).Methods("GET")
	s.AddApiRoute("/tasks/{id:[0-9]+}", handler.GetOneTaskHandler).Methods("GET")
	s.AddApiRoute("/tasks/{id:[0-9]+}", handler.UpdateTaskHandler).Methods("PUT")
	s.AddStatics("/", "static")
}
