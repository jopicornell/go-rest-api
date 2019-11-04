package tasks

import (
	"github.com/jopicornell/go-rest-api/internals/api/tasks/handlers"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s *server.Server) {
	handler := handlers.New(s)
	s.AddRoute("/tasks", handler.GetTasksHandler)
	s.AddRoute("/tasks/{id:[0-9]+}", handler.GetOneTaskHandler)
	return
}
