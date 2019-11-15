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

func (s *AppointmentHandler) GetAppointmentsHandler(request server.Request) {
	appointments, err := s.appointmentService.GetAppointments()
	if err != nil {
		log.Println(fmt.Errorf("error getting appointments: %w", err))
		request.Respond(http.StatusInternalServerError)
		return
	}
	request.RespondJSON(http.StatusOK, appointments)
}

func (s *AppointmentHandler) GetOneAppointmentHandler(request server.Request) {
	id := request.GetParamUInt("id")
	appointment, err := s.appointmentService.GetAppointment(uint(id))
	if err != nil {
		log.Println(fmt.Errorf("error getting appointment(%d): %w", id, err))
		request.Respond(http.StatusInternalServerError)
		return
	}
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
	appointment, err := s.appointmentService.UpdateAppointment(uint(id), appointment)
	if err != nil {
		log.Println(fmt.Errorf("error getting appointment(%d): %w", id, err))
		request.Respond(http.StatusInternalServerError)
		return
	}
	if appointment == nil {
		request.Respond(http.StatusNotFound)
		return
	}
	request.RespondJSON(http.StatusOK, appointment)
}

func (s *AppointmentHandler) CreateAppointmentHandler(request server.Request) {
	var appointment *models.Appointment
	request.GetBodyMarshalled(&appointment)
	if appointment, err := s.appointmentService.CreateAppointment(appointment, request.GetUser()); err == nil {
		request.RespondJSON(http.StatusCreated, appointment)
	} else {
		log.Println(fmt.Errorf("error creating appointment %+v: %w", appointment, err))
		request.Respond(http.StatusInternalServerError)
	}
}

func (s *AppointmentHandler) DeleteAppointmentHandler(request server.Request) {
	id := request.GetParamUInt("id")
	if err := s.appointmentService.DeleteAppointment(uint(id)); err == nil {
		request.Respond(http.StatusOK)
	} else {
		log.Println(fmt.Errorf("error deleting appointment %d: %w", id, err))
		request.Respond(http.StatusInternalServerError)
	}
}
