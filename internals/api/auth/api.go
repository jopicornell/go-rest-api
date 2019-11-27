package auth

import (
	"github.com/jopicornell/go-rest-api/internals/api/auth/handlers"
	server "github.com/jopicornell/go-rest-api/pkg/server"
)

func Start() {
	appServer := server.Initialize()
	defer appServer.Close()
	appServer.AddHandler(&handlers.AuthHandler{})
	appServer.ListenAndServe()
}
