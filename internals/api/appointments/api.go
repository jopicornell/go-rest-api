package appointments

import (
	"github.com/jopicornell/go-rest-api/internals/api/appointments/handlers"
	server "github.com/jopicornell/go-rest-api/pkg/server"
)

func Start() {
	appServer := server.Initialize()
	defer appServer.Close()
	appServer.AddHandler(&handlers.AppointmentHandler{})
	appServer.ListenAndServe()
}
