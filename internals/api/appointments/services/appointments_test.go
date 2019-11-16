package services

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	dbMock, _ := mockDB(t)
	appointmentsService := New(dbMock)
	if appointmentsService == nil {
		t.Errorf("New service should not be null")
	}
	reflectServiceDb := reflect.ValueOf(appointmentsService).Elem().FieldByName("db")
	if reflectServiceDb.IsNil() {
		t.Errorf("db field in service should not be nil")
	}
}

func TestAppointmentService_GetAppointments(t *testing.T) {
	t.Run("should panic if db throws", getAppointmentsShouldPanicIfDbThrows)
	t.Run("should return empty slice if no rows", getAppointmentsShouldReturnEmptySliceIfNoRows)
	t.Run("should return list of appointments if all went ok", getAppointmentsShouldReturnListOfAppointments)
}

func TestAppointmentService_GetAppointment(t *testing.T) {
	t.Run("should panic if db throws", getAppointmentShouldPanicIfDbThrows)
	t.Run("should return nil if no rows", getAppointmentShouldReturnNilIfNoRows)
	t.Run("should return appointment if all went ok", getAppointmentShouldReturnAppointment)
}

func TestAppointmentService_CreateAppointment(t *testing.T) {
	t.Run("should throw if db throws and rollback", createAppointmentShouldThrowIfDbThrows)
	t.Run("should throw if appointment to create is null", createAppointmentShouldThrowIfAppointmentIsNull)
	t.Run("should return created appointment and commit", createAppointmentShouldReturnAppointmentAndCommit)
}

func TestAppointmentService_UpdateAppointment(t *testing.T) {
	t.Run("should throw if db throws and rollback", updateAppointmentShouldThrowIfDbThrows)
	t.Run("should throw if appointment to updateis null", updateAppointmentShouldPanicIfAppointmentIsNull)
	t.Run("should return updated appointment and commit", updateAppointmentShouldReturnAppointmentAndCommit)
}

func TestAppointmentService_DeleteAppointment(t *testing.T) {
	t.Run("should throw if db throws and rollback", deleteAppointmentShouldThrowIfDbThrows)
	t.Run("should return no error and commit", deleteAppointmentShouldExecuteAndCommit)
}

func getAppointmentsShouldPanicIfDbThrows(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	expected := errors.New("test error")
	mock.ExpectQuery("SELECT \\* from appointments").WillReturnError(expected)
	defer func() {
		if recover := recover(); recover != nil {
			if recover != expected {
				t.Errorf("error expected %v got %+v", expected, recover)
			}
		} else {
			t.Errorf("Error should have been thrown")
		}
	}()
	result := appointmentService.GetAppointments()
	if result != nil {
		t.Errorf("result expected to be nil, found %+v", result)
	}
}

func getAppointmentsShouldReturnEmptySliceIfNoRows(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	mock.ExpectQuery("SELECT \\* from appointments").WillReturnRows(&sqlmock.Rows{})

	if got := appointmentService.GetAppointments(); got != nil {
		if len(got) != 0 {
			t.Errorf("expected result to be empty, got %+v", got)
		}
	} else {
		t.Errorf("result should not be empty")
	}
}

func getAppointmentsShouldReturnListOfAppointments(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	rows := buildAppointmentRows()
	expected := addAppointmentRows(rows, 5)
	mock.ExpectQuery("SELECT \\* from appointments").WillReturnRows(rows)

	if got := appointmentService.GetAppointments(); got != nil {
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected result to be %+v, got %+v", expected, got)
		}
	} else {
		t.Errorf("result should not be empty")
	}
}

func getAppointmentShouldPanicIfDbThrows(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	expected := errors.New("test error")
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from appointments").WithArgs(id).WillReturnError(expected)
	defer func() {
		if recover := recover(); recover != nil {
			if recover != expected {
				t.Errorf("error expected %v got %+v", expected, recover)
			}
		} else {
			t.Errorf("Error should have been thrown")
		}
	}()
	if result := appointmentService.GetAppointment(id); result != nil {
		if result != nil {
			t.Errorf("result expected to be nil got %+v", result)
		}
	}
}

func getAppointmentShouldReturnNilIfNoRows(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from appointments").WithArgs(id).WillReturnError(sql.ErrNoRows)

	if got := appointmentService.GetAppointment(id); got != nil {
		t.Errorf("result should be nill, got %+v", got)
	}
}

func getAppointmentShouldReturnAppointment(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	rows := buildAppointmentRows()
	expected := addAppointmentRows(rows, 1)[0]
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from appointments").WithArgs(id).WillReturnRows(rows)

	if got := appointmentService.GetAppointment(id); got != nil {
		if reflect.DeepEqual(got, expected) {
			t.Errorf("expected result to be %+v, got %+v", expected, got)
		}
	} else {
		t.Errorf("result should not be empty")
	}
}

func createAppointmentShouldReturnAppointmentAndCommit(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	appointment := createFakeAppointment()
	user := servertesting.CreateFakeUser()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO appointments (.*)").WithArgs(
		appointment.StartDate, appointment.Duration, appointment.EndDate, appointment.Status, user.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if got := appointmentService.CreateAppointment(appointment, user); got != nil {
		if appointment != got {
			t.Errorf("expected result to be %+v, got %+v", appointment, got)
		}
	} else {
		t.Errorf("expected response not to be null")
	}
}

func createAppointmentShouldThrowIfAppointmentIsNull(t *testing.T) {

}

func createAppointmentShouldThrowIfDbThrows(t *testing.T) {

}

func updateAppointmentShouldReturnAppointmentAndCommit(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	appointment := createFakeAppointment()
	id := uint(1)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE appointments SET (.*)").WithArgs(
		appointment.StartDate, appointment.Duration, appointment.EndDate, appointment.Status, appointment.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if got := appointmentService.UpdateAppointment(id, appointment); got != nil {
		if appointment != got {
			t.Errorf("expected result to be %+v, got %+v", appointment, got)
		}
	} else {
		t.Errorf("expected response not to be null")
	}
}

func updateAppointmentShouldThrowIfDbThrows(t *testing.T) {
	dbMock, mock := mockDB(t)
	appointmentService := New(dbMock)
	appointment := createFakeAppointment()
	id := uint(1)
	expectedError := errors.New("test error")
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE appointments SET (.*)").WithArgs(
		appointment.StartDate, appointment.Duration, appointment.EndDate, appointment.Status, appointment.ID,
	).WillReturnError(expectedError)
	mock.ExpectRollback()
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Errorf("Error should have been thrown")
		}
	}()
	if got := appointmentService.UpdateAppointment(id, appointment); got != nil {
		t.Errorf("result should have been nil")
	}
}

func updateAppointmentShouldPanicIfAppointmentIsNull(t *testing.T) {
	dbMock, _ := mockDB(t)
	appointmentService := New(dbMock)
	id := uint(1)
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Errorf("Error should have been thrown")
		}
	}()
	if got := appointmentService.UpdateAppointment(id, nil); got != nil {
		t.Errorf("expected updateAppointment to throw")
	}
}

func deleteAppointmentShouldThrowIfDbThrows(t *testing.T) {

}

func deleteAppointmentShouldExecuteAndCommit(t *testing.T) {

}

func mockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func buildAppointmentRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "start_date", "duration", "end_date", "status", "user_id"})
}

func addAppointmentRows(rows *sqlmock.Rows, numRows uint) []models.Appointment {
	var appointments []models.Appointment
	var appointment *models.Appointment
	for ; numRows > 0; numRows-- {
		appointment = createFakeAppointment()

		appointments = append(appointments, *appointment)
		rows.AddRow(appointment.ID, appointment.StartDate, appointment.Duration, appointment.EndDate, appointment.Status, appointment.UserId)
	}
	return appointments
}

func createFakeAppointment() *models.Appointment {
	return &models.Appointment{
		ID:        uint16(faker.UnixTime()),
		StartDate: time.Time{},
		Duration:  0,
		EndDate:   time.Time{},
		Status:    models.StatusPending,
		UserId:    0,
	}
}
