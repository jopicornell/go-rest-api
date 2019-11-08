package tasks

import (
	server "github.com/jopicornell/go-rest-api/pkg/server"
)

func Start() {
	appServer := server.Initialize()
	defer appServer.Close()
	ConfigureRoutes(appServer)
	appServer.ListenAndServe()
}
