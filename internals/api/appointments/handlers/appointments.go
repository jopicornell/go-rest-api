package handlers

import (
	"github.com/jopicornell/go-rest-api/internals/api/appointments/services"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

type AppointmentHandler struct {
	server.Handler
	appointmentService services.AppointmentsService
}

func New(s server.Server) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: services.New(s.GetRelationalDatabase()),
	}
}

func (s *AppointmentHandler) GetAppointmentsHandler(response server.Response, request server.Request) {
	if appointments, err := s.appointmentService.GetAppointments(); err == nil {
		response.RespondJSON(http.StatusOK, appointments)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (s *AppointmentHandler) GetOneAppointmentHandler(response server.Response, request server.Request) {
	id := request.GetParamUInt("id")
	if appointment, err := s.appointmentService.GetAppointment(uint(id)); err == nil {
		if appointment == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, appointment)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (s *AppointmentHandler) UpdateAppointmentHandler(response server.Response, request server.Request) {
	id := request.GetParamUInt("id")
	var appointment *models.Appointment
	request.GetBodyMarshalled(&appointment)
	if appointment, err := s.appointmentService.UpdateAppointment(uint(id), appointment); err == nil {
		if appointment == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, appointment)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError})
	}

}

func (s *AppointmentHandler) CreateAppointmentHandler(response server.Response, request server.Request) {
	var appointment *models.Appointment
	request.GetBodyMarshalled(&appointment)
	if appointment, err := s.appointmentService.CreateAppointment(appointment, request.GetUser()); err == nil {
		response.RespondJSON(http.StatusCreated, appointment)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}

}

func (s *AppointmentHandler) DeleteAppointmentHandler(response server.Response, request server.Request) {
	id := request.GetParamUInt("id")
	if err := s.appointmentService.DeleteAppointment(uint(id)); err == nil {
		response.Respond(http.StatusNoContent)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}
