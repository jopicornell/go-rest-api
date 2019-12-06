package handlers

import (
	"encoding/json"
	goErrors "errors"
	"github.com/bxcodec/faker/v3"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type PictureServiceMock struct {
	Picture      *models.Picture
	Pictures     []models.Picture
	errorToThrow error
}

func (ts *PictureServiceMock) UpdatePicture(id uint, Picture *models.Picture) (*models.Picture, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.Picture, nil
}

func (ts *PictureServiceMock) DeletePicture(id uint) error {
	if ts.errorToThrow != nil {
		return ts.errorToThrow
	}
	return nil
}

func (ts *PictureServiceMock) CreatePicture(Picture *models.Picture, user *models.User) (*models.Picture, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.Picture, nil
}

func (ts *PictureServiceMock) GetPicture(id uint) (*models.Picture, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.Picture, nil
}

func (ts *PictureServiceMock) GetPictures() ([]models.Picture, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.Pictures, nil
}

func TestPictureHandler_GetOnePictureHandler(t *testing.T) {
	t.Run("should throw error if service is missing", panicErrorServiceMissing)
	t.Run("should throw InternalServerError if service fails", internalErrorIfServiceFailsReturningPicture)
	t.Run("should throw a NotFound if service says so", notFoundIfServiceFoundsNothing)
	t.Run("should return Picture returned by the service", returnPictureByService)
}

func TestPictureHandler_GetPicturesHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", internalErrorIfServiceFailsReturningPictures)
	t.Run("should return Pictures returned by the service", returnPicturesByService)
}

func TestPictureHandler_UpdatePictureHandler(t *testing.T) {
	t.Run("should throw a NotFound if service says so", updatePictureShouldThrowNotFound)
	t.Run("should throw InternalServer if some error is raised by the service", updatePictureShouldThrowIfServiceFails)
	t.Run("should return Picture updated by the service", updatePictureShouldReturnUpdatedPicture)
}

func TestPictureHandler_CreatePictureHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", createPictureShouldThrowIfServiceFails)
	t.Run("should return Picture created by the service", createPictureShouldReturnCreatedPicture)
}

func TestPictureHandler_DeletePictureHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", deletePictureThrowIfError)
	t.Run("should return no content and no error", deleteSuccessAndReturnNoContent)
}

func panicErrorServiceMissing(t *testing.T) {
	handler := &PictureHandler{}
	response := servertesting.NewResponse()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/Pictures/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("function did not panic")
		}
	}()
	handler.GetOnePictureHandler(response, request)
}

func returnPictureByService(t *testing.T) {
	handler := &PictureHandler{}
	expected := &models.Picture{
		ID:        0,
		StartDate: time.Time{},
		Duration:  0,
		EndDate:   time.Time{},
		Status:    "",
		UserId:    0,
	}
	handler.pictureService = &PictureServiceMock{Picture: expected}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/Pictures/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetOnePictureHandler(recorder, context)
	var got *models.Picture
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

func returnPicturesByService(t *testing.T) {
	handler := &PictureHandler{}
	expected := []models.Picture{
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
	handler.pictureService = &PictureServiceMock{Pictures: expected}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/Pictures/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetPicturesHandler(recorder, context)
	var got []models.Picture
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed (%s): %w", recorder.Body.String(), err)
	}
	if len(got) != len(expected) {
		t.Errorf("expected length of result (%d) to be %d", len(got), len(expected))
	}

}

func internalErrorIfServiceFailsReturningPicture(t *testing.T) {
	handler := &PictureHandler{}
	var expected *models.Picture = nil
	handler.pictureService = &PictureServiceMock{
		Picture:      expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/Pictures/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetOnePictureHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}
}

func internalErrorIfServiceFailsReturningPictures(t *testing.T) {
	handler := &PictureHandler{}
	var expected []models.Picture = nil
	handler.pictureService = &PictureServiceMock{
		Pictures:     expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/Pictures/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetPicturesHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}
}

func notFoundIfServiceFoundsNothing(t *testing.T) {
	handler := &PictureHandler{}
	handler.pictureService = &PictureServiceMock{}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("GET", "/Pictures/1", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.GetOnePictureHandler(recorder, context)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %+v got %+v", http.StatusNotFound, recorder.Code)
	}
}

func updatePictureShouldReturnUpdatedPicture(t *testing.T) {
	handler := &PictureHandler{}
	expected := createFakePicture()
	handler.pictureService = &PictureServiceMock{
		Picture:      expected,
		errorToThrow: nil,
	}
	var PictureJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		PictureJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling Picture %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("PUT", "/Pictures/1", strings.NewReader(PictureJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.UpdatePictureHandler(recorder, context)
	var got *models.Picture
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

func updatePictureShouldThrowIfServiceFails(t *testing.T) {
	handler := &PictureHandler{}
	Picture := createFakePicture()
	errorToThrow := goErrors.New("test error")
	handler.pictureService = &PictureServiceMock{
		Picture:      Picture,
		errorToThrow: errorToThrow,
	}
	var PictureJSON string
	if marshallResult, err := json.Marshal(Picture); err == nil {
		PictureJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling Picture %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("PUT", "/Pictures/1", strings.NewReader(PictureJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.UpdatePictureHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}

}

func updatePictureShouldThrowNotFound(t *testing.T) {
	handler := &PictureHandler{}
	Picture := createFakePicture()
	handler.pictureService = &PictureServiceMock{
		Picture:      nil,
		errorToThrow: nil,
	}
	var PictureJSON string
	if marshallResult, err := json.Marshal(Picture); err == nil {
		PictureJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling Picture %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("PUT", "/Pictures/1", strings.NewReader(PictureJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.UpdatePictureHandler(recorder, context)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %+v got %+v", http.StatusNotFound, recorder.Code)
	}
}

func createPictureShouldReturnCreatedPicture(t *testing.T) {
	handler := &PictureHandler{}
	expected := createFakePicture()
	handler.pictureService = &PictureServiceMock{
		Picture: expected,
	}
	var PictureJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		PictureJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling Picture %w", err)
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("POST", "/Pictures", strings.NewReader(PictureJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.CreatePictureHandler(recorder, context)
	var got *models.Picture
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

func createPictureShouldThrowIfServiceFails(t *testing.T) {
	handler := &PictureHandler{}

	Picture := createFakePicture()
	var PictureJSON string
	if marshallResult, err := json.Marshal(Picture); err == nil {
		PictureJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling Picture %w", err)
	}
	handler.pictureService = &PictureServiceMock{
		errorToThrow: goErrors.New("test error"),
	}
	recorder := servertesting.NewResponse()
	context := servertesting.NewRequest(
		httptest.NewRequest("POST", "/Pictures", strings.NewReader(PictureJSON)),
		map[string]string{
			"id": "1",
		},
	)
	handler.CreatePictureHandler(recorder, context)
	if err := recorder.ErrorCalled(); err != nil {
		if err.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected status code to be %d returned %d", http.StatusInternalServerError, err.StatusCode)
		}
	} else {
		t.Errorf("expected error to be called")
	}
}

func deleteSuccessAndReturnNoContent(t *testing.T) {
	handler := &PictureHandler{}
	handler.pictureService = &PictureServiceMock{}
	recorder := servertesting.NewResponse()
	request := servertesting.NewRequest(
		httptest.NewRequest("DELETE", "/Pictures", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.DeletePictureHandler(recorder, request)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected status code to be %d got %d", http.StatusNoContent, recorder.Code)
	}
}

func deletePictureThrowIfError(t *testing.T) {
	handler := &PictureHandler{}
	handler.pictureService = &PictureServiceMock{}
	recorder := servertesting.NewResponse()
	request := servertesting.NewRequest(
		httptest.NewRequest("DELETE", "/Pictures", nil),
		map[string]string{
			"id": "1",
		},
	)
	handler.DeletePictureHandler(recorder, request)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusNoContent {
		t.Errorf("expected status code to be %d got %d", http.StatusNoContent, recorder.Code)
	}
}

func createFakePicture() *models.Picture {
	return &models.Picture{
		ID:        uint16(faker.UnixTime()),
		StartDate: time.Time{},
		Duration:  0,
		EndDate:   time.Time{},
		Status:    models.StatusPending,
		UserId:    0,
	}
}
