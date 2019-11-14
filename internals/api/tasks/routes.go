package tasks

import (
	"github.com/jopicornell/go-rest-api/internals/api/tasks/handlers"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s server.Server) {
	handler := handlers.New(s)
	group := s.GetRouter().AddGroup("/api")
	group.AddRoute("/tasks", handler.GetTasksHandler).Methods("GET")
	group.AddRoute("/tasks", handler.CreateTaskHandler).Methods("POST")
	group.AddRoute("/tasks/{id:[0-9]+}", handler.DeleteTaskHandler).Methods("DELETE")
	group.AddRoute("/tasks/{id:[0-9]+}", handler.GetOneTaskHandler).Methods("GET")
	group.AddRoute("/tasks/{id:[0-9]+}", handler.UpdateTaskHandler).Methods("PUT")
	s.AddStatics("/", "static")
}
