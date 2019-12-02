package handlers

import (
	"github.com/jopicornell/go-rest-api/internals/api/middlewares"
	"github.com/jopicornell/go-rest-api/internals/api/services"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

type AppointmentHandler struct {
	server             server.Server
	appointmentService services.AppointmentsService
}

func (a *AppointmentHandler) Initialize(server server.Server) {
	a.server = server
	a.appointmentService = services.NewAppointmentService(server.GetRelationalDatabase())
}

func (a *AppointmentHandler) ConfigureRoutes() server.Router {
	appointments := server.NewRouter().AddGroup("/appointments")
	appointments.Use(&middlewares.UserMiddleware{}, &middlewares.UserMiddleware{})
	appointments.AddRoute("", a.GetAppointmentsHandler).Methods("GET")
	appointments.AddRoute("", a.CreateAppointmentHandler).Methods("POST")
	appointments.AddRoute("/{id:[0-9]+}", a.DeleteAppointmentHandler).Methods("DELETE")
	appointments.AddRoute("/{id:[0-9]+}", a.GetOneAppointmentHandler).Methods("GET")
	appointments.AddRoute("/{id:[0-9]+}", a.UpdateAppointmentHandler).Methods("PUT")
	return appointments
}

func (a *AppointmentHandler) GetAppointmentsHandler(response server.Response, request server.Request) {
	if appointments, err := a.appointmentService.GetAppointments(); err == nil {
		response.RespondJSON(http.StatusOK, appointments)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *AppointmentHandler) GetOneAppointmentHandler(response server.Response, request server.Request) {
	id := request.GetParamUInt("id")
	if appointment, err := a.appointmentService.GetAppointment(uint(id)); err == nil {
		if appointment == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, appointment)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *AppointmentHandler) UpdateAppointmentHandler(response server.Response, request server.Request) {
	id := request.GetParamUInt("id")
	var appointment *models.Appointment
	request.GetBodyMarshalled(&appointment)
	if appointment, err := a.appointmentService.UpdateAppointment(uint(id), appointment); err == nil {
		if appointment == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, appointment)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError})
	}

}

func (a *AppointmentHandler) CreateAppointmentHandler(response server.Response, request server.Request) {
	var appointment *models.Appointment
	request.GetBodyMarshalled(&appointment)
	if appointment, err := a.appointmentService.CreateAppointment(appointment, request.GetUser()); err == nil {
		response.RespondJSON(http.StatusCreated, appointment)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}

}

func (a *AppointmentHandler) DeleteAppointmentHandler(response server.Response, request server.Request) {
	id := request.GetParamUInt("id")
	if err := a.appointmentService.DeleteAppointment(uint(id)); err == nil {
		response.Respond(http.StatusNoContent)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}
