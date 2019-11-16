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

func (s *AppointmentHandler) GetAppointmentsHandler(request server.Request) {
	appointments := s.appointmentService.GetAppointments()
	request.RespondJSON(http.StatusOK, appointments)
}

func (s *AppointmentHandler) GetOneAppointmentHandler(request server.Request) {
	id := request.GetParamUInt("id")
	appointment := s.appointmentService.GetAppointment(uint(id))
	if appointment == nil {
		request.Respond(http.StatusNotFound)
		return
	}
	request.RespondJSON(http.StatusOK, appointment)
}

func (s *AppointmentHandler) UpdateAppointmentHandler(request server.Request) {
	id := request.GetParamUInt("id")
	var appointment *models.Appointment
	request.GetBodyMarshalled(&appointment)
	appointment = s.appointmentService.UpdateAppointment(uint(id), appointment)
	if appointment == nil {
		request.Respond(http.StatusNotFound)
		return
	}
	request.RespondJSON(http.StatusOK, appointment)
}

func (s *AppointmentHandler) CreateAppointmentHandler(request server.Request) {
	var appointment *models.Appointment
	request.GetBodyMarshalled(&appointment)
	appointment = s.appointmentService.CreateAppointment(appointment, request.GetUser())
	request.RespondJSON(http.StatusCreated, appointment)

}

func (s *AppointmentHandler) DeleteAppointmentHandler(request server.Request) {
	id := request.GetParamUInt("id")
	s.appointmentService.DeleteAppointment(uint(id))
	request.Respond(http.StatusNoContent)
}
