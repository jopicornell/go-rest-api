package handlers

import (
	"encoding/json"
	goErrors "errors"
	"github.com/bxcodec/faker/v3"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type AppointmentHandlerMock struct {
	AppointmentHandler
}

type AppointmentServiceMock struct {
	appointment  *models.Appointment
	appointments []models.Appointment
	errorToThrow error
}

func (ts *AppointmentServiceMock) UpdateAppointment(id uint, appointment *models.Appointment) (*models.Appointment, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.appointment, nil
}

func (ts *AppointmentServiceMock) DeleteAppointment(id uint) error {
	if ts.errorToThrow != nil {
		return ts.errorToThrow
	}
	return nil
}

func (ts *AppointmentServiceMock) CreateAppointment(appointment *models.Appointment, user *models.User) (*models.Appointment, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.appointment, nil
}

func (ts *AppointmentServiceMock) GetAppointment(id uint) (*models.Appointment, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.appointment, nil
}

func (ts *AppointmentServiceMock) GetAppointments() ([]models.Appointment, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.appointments, nil
}

func TestAppointmentHandler_New(t *testing.T) {
	t.Run("should construct a new AppointmentHandler given the server", shouldReturnConstructedHandler)
}

func TestAppointmentHandler_GetOneAppointmentHandler(t *testing.T) {
	t.Run("should throw error if service is missing", panicErrorServiceMissing)
	t.Run("should throw InternalServerError if service fails", internalErrorIfServiceFailsReturningAppointment)
	t.Run("should throw a NotFound if service says so", notFoundIfServiceFoundsNothing)
	t.Run("should return appointment returned by the service", returnAppointmentByService)
}

func TestAppointmentHandler_GetAppointmentsHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", internalErrorIfServiceFailsReturningAppointments)
	t.Run("should return appointments returned by the service", returnAppointmentsByService)
}

func TestAppointmentHandler_UpdateAppointmentHandler(t *testing.T) {
	t.Run("should throw a NotFound if service says so", updateAppointmentShouldThrowNotFound)
	t.Run("should throw InternalServer if some error is raised by the service", updateAppointmentShouldThrowIfServiceFails)
	t.Run("should return appointment updated by the service", updateAppointmentShouldReturnUpdatedAppointment)
}

func TestAppointmentHandler_CreateAppointmentHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", createAppointmentShouldThrowIfServiceFails)
	t.Run("should return appointment created by the service", createAppointmentShouldReturnCreatedAppointment)
}

func panicErrorServiceMissing(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/appointments/1", nil),
		httptest.NewRecorder(),
		map[string]string{
			"id": "1",
		},
	)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("function did not panic")
		}
	}()
	mock.GetOneAppointmentHandler(context)
}

func returnAppointmentByService(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	expected := &models.Appointment{
		ID:        0,
		StartDate: time.Time{},
		Duration:  0,
		EndDate:   time.Time{},
		Status:    "",
		UserId:    0,
	}
	mock.appointmentService = &AppointmentServiceMock{appointment: expected}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/appointments/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetOneAppointmentHandler(context)
	var got *models.Appointment
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed: %w %s", err, recorder.Body.String())
	}
	if !reflect.DeepEqual(*got, *expected) {
		t.Errorf("expected %+v got %+v", *expected, *got)
	}
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d got %d", recorder.Code, http.StatusOK)
	}
}

func returnAppointmentsByService(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	expected := []models.Appointment{
		{
			ID:        0,
			StartDate: time.Time{},
			Duration:  0,
			EndDate:   time.Time{},
			Status:    "",
			UserId:    0,
		},
		{
			ID:        0,
			StartDate: time.Time{},
			Duration:  0,
			EndDate:   time.Time{},
			Status:    "",
			UserId:    0,
		},
	}
	mock.appointmentService = &AppointmentServiceMock{appointments: expected}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/appointments/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetAppointmentsHandler(context)
	var got []models.Appointment
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed (%s): %w", recorder.Body.String(), err)
	}
	if len(got) != len(expected) {
		t.Errorf("expected length of result (%d) to be %d", len(got), len(expected))
	}

}

func internalErrorIfServiceFailsReturningAppointment(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	var expected *models.Appointment = nil
	mock.appointmentService = &AppointmentServiceMock{
		appointment:  expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/appointments/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetOneAppointmentHandler(context)
	if recorder.Code != 500 {
		t.Errorf("expected status code to be 500 got %+v", recorder.Code)
	}
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v (%s)", recorder.Body.Len(), recorder.Body.String())
	}
}

func internalErrorIfServiceFailsReturningAppointments(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	var expected []models.Appointment = nil
	mock.appointmentService = &AppointmentServiceMock{
		appointments: expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/appointments/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetAppointmentsHandler(context)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status code to be %d got %d", http.StatusInternalServerError, recorder.Code)
	}
}

func notFoundIfServiceFoundsNothing(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	mock.appointmentService = &AppointmentServiceMock{}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/appointments/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetOneAppointmentHandler(context)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %+v got %+v", http.StatusNotFound, recorder.Code)
	}
}

func shouldReturnConstructedHandler(t *testing.T) {
	serverMock := &servertesting.ServerMock{
		Config: config.Config{},
	}
	appointmentHandler := New(serverMock)
	if appointmentHandler == nil {
		t.Errorf("appointment handler should not be null")
		return
	}
	if appointmentHandler.appointmentService == nil {
		t.Errorf("appointment handler created without the service")
	}
}

func updateAppointmentShouldReturnUpdatedAppointment(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	expected := createFakeAppointment()
	mock.appointmentService = &AppointmentServiceMock{
		appointment:  expected,
		errorToThrow: nil,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("PUT", "/appointments/1", strings.NewReader(appointmentJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.UpdateAppointmentHandler(context)
	var got *models.Appointment
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed: %w", err)
	}
	if !reflect.DeepEqual(*got, *expected) {
		t.Errorf("expected %+v got %+v", expected, got)
	}
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d got %d", recorder.Code, http.StatusOK)
	}

}

func updateAppointmentShouldThrowIfServiceFails(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	appointment := createFakeAppointment()
	errorToThrow := goErrors.New("test error")
	mock.appointmentService = &AppointmentServiceMock{
		appointment:  appointment,
		errorToThrow: errorToThrow,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(appointment); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("PUT", "/appointments/1", strings.NewReader(appointmentJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.UpdateAppointmentHandler(context)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status code to be %d got %d", http.StatusInternalServerError, recorder.Code)
	}

}

func updateAppointmentShouldThrowNotFound(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	appointment := createFakeAppointment()
	mock.appointmentService = &AppointmentServiceMock{
		appointment:  nil,
		errorToThrow: nil,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(appointment); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("PUT", "/appointments/1", strings.NewReader(appointmentJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.UpdateAppointmentHandler(context)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %+v got %+v", http.StatusNotFound, recorder.Code)
	}
}

func createAppointmentShouldReturnCreatedAppointment(t *testing.T) {
	mock := &AppointmentHandlerMock{}
	expected := createFakeAppointment()
	mock.appointmentService = &AppointmentServiceMock{
		appointment:  expected,
		errorToThrow: nil,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("POST", "/appointments", strings.NewReader(appointmentJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.CreateAppointmentHandler(context)
	var got *models.Appointment
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed: %w", err)
		return
	}
	if *got != *expected {
		t.Errorf("expected %+v got %+v", expected, got)
	}
	if recorder.Code != http.StatusCreated {
		t.Errorf("expected status code %d got %d", http.StatusCreated, recorder.Code)
	}
}

func createAppointmentShouldThrowIfServiceFails(t *testing.T) {
	mock := &AppointmentHandlerMock{}

	appointment := createFakeAppointment()
	var appointmentJSON string
	if marshallResult, err := json.Marshal(appointment); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	mock.appointmentService = &AppointmentServiceMock{
		appointment:  nil,
		appointments: nil,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("POST", "/appointments", strings.NewReader(appointmentJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.CreateAppointmentHandler(context)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status code to be %d got %d", http.StatusInternalServerError, recorder.Code)
	}
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
