package services

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/models"
)

type AppointmentsService interface {
	GetAppointments() ([]models.Appointment, error)
	GetAppointment(uint) (*models.Appointment, error)
	UpdateAppointment(uint, *models.Appointment) (*models.Appointment, error)
	CreateAppointment(*models.Appointment, *models.User) (*models.Appointment, error)
	DeleteAppointment(uint) error
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

func (s *appointmentService) GetAppointments() (appointments []models.Appointment, err error) {
	appointments = []models.Appointment{}
	if err = s.db.Select(&appointments, "SELECT * from appointments"); err != nil {
		return nil, err
	}
	return appointments, nil
}

func (s *appointmentService) GetAppointment(id uint) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := s.db.Get(&appointment, "SELECT * from appointments where id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &appointment, nil
}

func (s *appointmentService) CreateAppointment(appointment *models.Appointment, user *models.User) (*models.Appointment, error) {
	if appointment == nil {
		return nil, AppointmentNullError
	}
	insertQuery := "INSERT INTO appointments (id, start_date, duration, end_date, status, user_id) VALUES (?, ?, ?, 0)"
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(insertQuery, appointment.StartDate, appointment.Duration, appointment.EndDate, models.StatusPending, user.ID); err == nil {
		err = tx.Commit()
		return appointment, err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, err
		}
		return nil, err
	}
}

func (s *appointmentService) UpdateAppointment(id uint, appointment *models.Appointment) (*models.Appointment, error) {
	if appointment == nil {
		return nil, AppointmentNullError
	}
	updateQuery := "UPDATE appointments SET start_date=?, duration=?, end_date=?, status=? where id = ?"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(updateQuery, appointment.StartDate, appointment.Duration, appointment.EndDate, appointment.ID); err == nil {
		err = tx.Commit()
		return appointment, err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, err
		}
		return nil, err
	}
}

func (s *appointmentService) DeleteAppointment(id uint) (err error) {
	deleteQuery := "DELETE FROM appointments WHERE id = ?"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(deleteQuery, id); err == nil {
		err = tx.Commit()
		return err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}
}
