package appointments

import (
	"github.com/jopicornell/go-rest-api/internals/api/appointments/handlers"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s server.Server) {
	handler := handlers.New(s)
	group := s.GetRouter().AddGroup("/api")
	group.AddRoute("/appointments", handler.GetAppointmentsHandler).Methods("GET")
	group.AddRoute("/appointments", handler.CreateAppointmentHandler).Methods("POST")
	group.AddRoute("/appointments/{id:[0-9]+}", handler.DeleteAppointmentHandler).Methods("DELETE")
	group.AddRoute("/appointments/{id:[0-9]+}", handler.GetOneAppointmentHandler).Methods("GET")
	group.AddRoute("/appointments/{id:[0-9]+}", handler.UpdateAppointmentHandler).Methods("PUT")
	s.AddStatics("/", "static")
}
