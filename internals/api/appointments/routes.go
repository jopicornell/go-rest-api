package appointments

import (
	"github.com/jopicornell/go-rest-api/internals/api/appointments/handlers"
	"github.com/jopicornell/go-rest-api/internals/middlewares"
	commonMiddlewares "github.com/jopicornell/go-rest-api/pkg/middlewares"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func ConfigureRoutes(s server.Server) {
	handler := handlers.New(s)
	privateGroup := s.GetRouter().AddGroup("")
	privateGroup.Use(commonMiddlewares.JWTPayload(s), middlewares.UserMiddleware(s))
	privateGroup.AddRoute("/appointments", handler.GetAppointmentsHandler).Methods("GET")
	privateGroup.AddRoute("/appointments", handler.CreateAppointmentHandler).Methods("POST")
	privateGroup.AddRoute("/appointments/{id:[0-9]+}", handler.DeleteAppointmentHandler).Methods("DELETE")
	privateGroup.AddRoute("/appointments/{id:[0-9]+}", handler.GetOneAppointmentHandler).Methods("GET")
	privateGroup.AddRoute("/appointments/{id:[0-9]+}", handler.UpdateAppointmentHandler).Methods("PUT")
	s.AddStatics("/", "static")
}
