package services

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/models"
)

type AppointmentsService interface {
	GetAppointments() []models.Appointment
	GetAppointment(uint) *models.Appointment
	UpdateAppointment(uint, *models.Appointment) *models.Appointment
	CreateAppointment(*models.Appointment, *models.User) *models.Appointment
	DeleteAppointment(uint)
}

type appointmentService struct {
	db *sqlx.DB
}

var AppointmentNullError = errors.New("appointment should not be null")

func New(db *sqlx.DB) AppointmentsService {
	return &appointmentService{
		db: db,
	}
}

func (s *appointmentService) GetAppointments() (appointments []models.Appointment) {
	appointments = []models.Appointment{}
	if err := s.db.Select(&appointments, "SELECT * from appointments"); err != nil {
		panic(err)
	}
	return appointments
}

func (s *appointmentService) GetAppointment(id uint) *models.Appointment {
	var appointment models.Appointment
	if err := s.db.Get(&appointment, "SELECT * from appointments where id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return &appointment
}

func (s *appointmentService) CreateAppointment(appointment *models.Appointment, user *models.User) *models.Appointment {
	insertQuery := "INSERT INTO appointments (start_date, duration, end_date, status, user_id) VALUES (?, ?, ?, ?, ?)"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(insertQuery, appointment.StartDate, appointment.Duration, appointment.EndDate, models.StatusPending, user.ID); err == nil {
		if err = tx.Commit(); err != nil {
			if err = tx.Rollback(); err != nil {
				panic(err)
			}
			panic(err)
		}
		return appointment
	} else {
		if errRollback := tx.Rollback(); errRollback != nil {

			panic(errRollback)
		} else {
			panic(err)
		}
	}
}

func (s *appointmentService) UpdateAppointment(id uint, appointment *models.Appointment) *models.Appointment {
	updateQuery := "UPDATE appointments SET start_date=?, duration=?, end_date=?, status=? where id = ?"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(updateQuery, appointment.StartDate, appointment.Duration, appointment.EndDate, appointment.Status, appointment.ID); err == nil {
		if err = tx.Commit(); err != nil {
			panic(err)
		}
		return appointment
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			panic(err)
		}
		panic(err)
	}
}

func (s *appointmentService) DeleteAppointment(id uint) {
	deleteQuery := "DELETE FROM appointments WHERE id = ?"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(deleteQuery, id); err == nil {
		if err = tx.Commit(); err != nil {
			panic(err)
		}
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			panic(errRollback)
		}
		panic(err)
	}
}
