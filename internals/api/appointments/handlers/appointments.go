package handlers

import (
	"fmt"
	"github.com/jopicornell/go-rest-api/internals/api/appointments/services"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"log"
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

func (s *AppointmentHandler) GetAppointmentsHandler(context server.Context) {
	appointments, err := s.appointmentService.GetAppointments()
	if err != nil {
		log.Println(fmt.Errorf("error getting appointments: %w", err))
		context.Respond(http.StatusInternalServerError)
		return
	}
	context.RespondJSON(http.StatusOK, appointments)
}

func (s *AppointmentHandler) GetOneAppointmentHandler(context server.Context) {
	id := context.GetParamUInt("id")
	appointment, err := s.appointmentService.GetAppointment(uint(id))
	if err != nil {
		log.Println(fmt.Errorf("error getting appointment(%d): %w", id, err))
		context.Respond(http.StatusInternalServerError)
		return
	}
	if appointment == nil {
		context.Respond(http.StatusNotFound)
		return
	}
	context.RespondJSON(http.StatusOK, appointment)
}

func (s *AppointmentHandler) UpdateAppointmentHandler(context server.Context) {
	id := context.GetParamUInt("id")
	var appointment *models.Appointment
	context.GetBodyMarshalled(&appointment)
	appointment, err := s.appointmentService.UpdateAppointment(uint(id), appointment)
	if err != nil {
		log.Println(fmt.Errorf("error getting appointment(%d): %w", id, err))
		context.Respond(http.StatusInternalServerError)
		return
	}
	if appointment == nil {
		context.Respond(http.StatusNotFound)
		return
	}
	context.RespondJSON(http.StatusOK, appointment)
}

func (s *AppointmentHandler) CreateAppointmentHandler(context server.Context) {
	var appointment *models.Appointment
	context.GetBodyMarshalled(&appointment)
	if appointment, err := s.appointmentService.CreateAppointment(appointment, context.GetUser()); err == nil {
		context.RespondJSON(http.StatusCreated, appointment)
	} else {
		log.Println(fmt.Errorf("error creating appointment %+v: %w", appointment, err))
		context.Respond(http.StatusInternalServerError)
	}
}

func (s *AppointmentHandler) DeleteAppointmentHandler(context server.Context) {
	id := context.GetParamUInt("id")
	if err := s.appointmentService.DeleteAppointment(uint(id)); err == nil {
		context.Respond(http.StatusOK)
	} else {
		log.Println(fmt.Errorf("error deleting appointment %d: %w", id, err))
		context.Respond(http.StatusInternalServerError)
	}
}
