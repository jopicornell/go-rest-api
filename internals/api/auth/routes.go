package tasks

import (
	"github.com/jopicornell/go-rest-api/internals/api/tasks/handlers"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s server.Server) {
	handler := handlers.New(s)
	s.AddApiRoute("/login", handler.CreateTaskHandler).Methods("POST")
	s.AddStatics("/", "static")
}
