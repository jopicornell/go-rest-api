package auth

import (
	"github.com/jopicornell/go-rest-api/internals/api/auth/handlers"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s server.Server) {
	handler := handlers.New(s)
	s.AddApiRoute("/login", handler.Login).Methods("POST")
}
