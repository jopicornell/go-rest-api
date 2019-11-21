package appointments

import (
	server "github.com/jopicornell/go-rest-api/pkg/server"
)

func Start() {
	appServer := server.Initialize()
	defer appServer.Close()
	appServer.ListenAndServe()
}
