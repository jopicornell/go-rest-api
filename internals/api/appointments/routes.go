package appointments

import (
	"github.com/jopicornell/go-rest-api/internals/api/appointments/handlers"
	"github.com/jopicornell/go-rest-api/internals/middlewares"
	commonMiddlewares "github.com/jopicornell/go-rest-api/pkg/middlewares"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s server.Server) {
	handler := handlers.New(s)
	appointments := s.GetRouter().AddGroup("/appointments")
	appointments.Use(commonMiddlewares.JWTPayload(s), middlewares.UserMiddleware(s))
	appointments.AddRoute("", handler.GetAppointmentsHandler).Methods("GET")
	appointments.AddRoute("", handler.CreateAppointmentHandler).Methods("POST")
	appointments.AddRoute("/{id:[0-9]+}", handler.DeleteAppointmentHandler).Methods("DELETE")
	appointments.AddRoute("/{id:[0-9]+}", handler.GetOneAppointmentHandler).Methods("GET")
	appointments.AddRoute("/{id:[0-9]+}", handler.UpdateAppointmentHandler).Methods("PUT")
	s.AddStatics("/", "static")
}
