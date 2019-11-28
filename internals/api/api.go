package api

import (
	"github.com/jopicornell/go-rest-api/internals/api/handlers"
	server "github.com/jopicornell/go-rest-api/pkg/server"
)

func Start() {
	appServer := server.Initialize()
	defer appServer.Close()
	appServer.AddHandler(&handlers.AppointmentHandler{})
	appServer.AddHandler(&handlers.AuthHandler{})
	appServer.ListenAndServe()
}