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

func TestAppointmentHandler_DeleteAppointmentHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", deleteAppointmentThrowIfError)
	t.Run("should return no content and no error", deleteSuccessAndReturnNoContent)
}

func panicErrorServiceMissing(t *testing.T) {
	handler := &AppointmentHandler{}
	response := servertesting.NewResponse()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/appointments/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("function did not panic")
		}
	}()
	handler.GetOneAppointmentHandler(response, request)
}

func returnAppointmentByService(t *testing.T) {
	handler := &AppointmentHandler{}
	expected := &models.Appointment{
		ID:        0,
		StartDate: time.Time{},
		Duration:  0,
		EndDate:   time.Time{},
		Status:    "",
		UserId:    0,
	}
	handler.appointmentService = &AppointmentServiceMock{appointment: expected}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/appointments/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetOneAppointmentHandler(recorder, context)
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
	handler := &AppointmentHandler{}
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
	handler.appointmentService = &AppointmentServiceMock{appointments: expected}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/appointments/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetAppointmentsHandler(recorder, context)
	var got []models.Appointment
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed (%s): %w", recorder.Body.String(), err)
	}
	if len(got) != len(expected) {
		t.Errorf("expected length of result (%d) to be %d", len(got), len(expected))
	}

}

func internalErrorIfServiceFailsReturningAppointment(t *testing.T) {
	handler := &AppointmentHandler{}
	var expected *models.Appointment = nil
	handler.appointmentService = &AppointmentServiceMock{
		appointment:  expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/appointments/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetOneAppointmentHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}
}

func internalErrorIfServiceFailsReturningAppointments(t *testing.T) {
	handler := &AppointmentHandler{}
	var expected []models.Appointment = nil
	handler.appointmentService = &AppointmentServiceMock{
		appointments: expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/appointments/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetAppointmentsHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}
}

func notFoundIfServiceFoundsNothing(t *testing.T) {
	handler := &AppointmentHandler{}
	handler.appointmentService = &AppointmentServiceMock{}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/appointments/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetOneAppointmentHandler(recorder, context)
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
	handler := &AppointmentHandler{}
	expected := createFakeAppointment()
	handler.appointmentService = &AppointmentServiceMock{
		appointment:  expected,
		errorToThrow: nil,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("PUT", "/appointments/1", strings.NewReader(appointmentJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.UpdateAppointmentHandler(recorder, context)
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
	handler := &AppointmentHandler{}
	appointment := createFakeAppointment()
	errorToThrow := goErrors.New("test error")
	handler.appointmentService = &AppointmentServiceMock{
		appointment:  appointment,
		errorToThrow: errorToThrow,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(appointment); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("PUT", "/appointments/1", strings.NewReader(appointmentJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.UpdateAppointmentHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}

}

func updateAppointmentShouldThrowNotFound(t *testing.T) {
	handler := &AppointmentHandler{}
	appointment := createFakeAppointment()
	handler.appointmentService = &AppointmentServiceMock{
		appointment:  nil,
		errorToThrow: nil,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(appointment); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("PUT", "/appointments/1", strings.NewReader(appointmentJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.UpdateAppointmentHandler(recorder, context)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %+v got %+v", http.StatusNotFound, recorder.Code)
	}
}

func createAppointmentShouldReturnCreatedAppointment(t *testing.T) {
	handler := &AppointmentHandler{}
	expected := createFakeAppointment()
	handler.appointmentService = &AppointmentServiceMock{
		appointment: expected,
	}
	var appointmentJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("POST", "/appointments", strings.NewReader(appointmentJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.CreateAppointmentHandler(recorder, context)
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
	handler := &AppointmentHandler{}

	appointment := createFakeAppointment()
	var appointmentJSON string
	if marshallResult, err := json.Marshal(appointment); err == nil {
		appointmentJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling appointment %w", err)
	}
	handler.appointmentService = &AppointmentServiceMock{
		errorToThrow: goErrors.New("test error"),
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("POST", "/appointments", strings.NewReader(appointmentJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.CreateAppointmentHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}
}

func deleteSuccessAndReturnNoContent(t *testing.T) {
	handler := &AppointmentHandler{}
	handler.appointmentService = &AppointmentServiceMock{}
	recorder := servertesting.NewResponse()
	request := servertesting.NewRequest(
		httptest.NewRequest("DELETE", "/appointments", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.DeleteAppointmentHandler(recorder, request)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected status code to be %d got %d", http.StatusNoContent, recorder.Code)
	}
}

func deleteAppointmentThrowIfError(t *testing.T) {
	handler := &AppointmentHandler{}
	handler.appointmentService = &AppointmentServiceMock{}
	recorder := servertesting.NewResponse()
	request := servertesting.NewRequest(
		httptest.NewRequest("DELETE", "/appointments", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.DeleteAppointmentHandler(recorder, request)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected status code to be %d got %d", http.StatusNoContent, recorder.Code)
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
