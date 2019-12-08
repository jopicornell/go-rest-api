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

func TestNewPictureService(t *testing.T) {
	dbMock, _ := mockPicturesDB(t)
	picturesService := NewPictureService(dbMock)
	if picturesService == nil {
		t.Errorf("New service should not be null")
	}
	reflectServiceDb := reflect.ValueOf(picturesService).Elem().FieldByName("db")
	if reflectServiceDb.IsNil() {
		t.Errorf("db field in service should not be nil")
	}
}

func TestPictureService_GetPictures(t *testing.T) {
	t.Run("should throw if db throws", getPicturesShouldThrowIfDbThrows)
	t.Run("should return empty slice if no rows", getPicturesShouldReturnEmptySliceIfNoRows)
	t.Run("should return list of pictures if all went ok", getPicturesShouldReturnListOfPictures)
}

func TestPictureService_GetPicture(t *testing.T) {
	t.Run("should throw if db throws", getPictureShouldThrowIfDbThrows)
	t.Run("should return nil if no rows", getPictureShouldReturnNilIfNoRows)
	t.Run("should return picture if all went ok", getPictureShouldReturnPicture)
}

func TestPictureService_CreatePicture(t *testing.T) {
	t.Run("should throw if db throws and rollback", createPictureShouldThrowIfDbThrows)
	t.Run("should throw if picture to create is null", createPictureShouldThrowIfPictureIsNull)
	t.Run("should return created picture and commit", createPictureShouldReturnPictureAndCommit)
}

func TestPictureService_UpdatePicture(t *testing.T) {
	t.Run("should throw if db throws and rollback", updatePictureShouldThrowIfDbThrows)
	t.Run("should throw if picture to updateis null", updatePictureShouldThrowIfPictureIsNull)
	t.Run("should return updated picture and commit", updatePictureShouldReturnPictureAndCommit)
}

func TestPictureService_DeletePicture(t *testing.T) {
	t.Run("should throw if db throws and rollback", deletePictureShouldThrowIfDbThrows)
	t.Run("should return no error and commit", deletePictureShouldExecuteAndCommit)
}

func getPicturesShouldThrowIfDbThrows(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	expected := errors.New("test error")
	mock.ExpectQuery("SELECT \\* from pictures").WillReturnError(expected)

	if _, got := pictureService.GetPictures(); got != nil {
		if got != expected {
			t.Errorf("error expected %v got %+v", expected, got)
		}
	} else {
		t.Errorf("Error should have been thrown")
	}
}

func getPicturesShouldReturnEmptySliceIfNoRows(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	mock.ExpectQuery("SELECT \\* from pictures").WillReturnRows(&sqlmock.Rows{})

	if got, err := pictureService.GetPictures(); got != nil {
		if err != nil {
			t.Errorf("expected err %+v to be nil", err)
		}
		if len(got) != 0 {
			t.Errorf("expected result to be empty, got %+v", got)
		}
	} else {
		t.Errorf("result should not be empty")
	}
}

func getPicturesShouldReturnListOfPictures(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	rows := buildPictureRows()
	expected := addPictureRows(rows, 5)
	mock.ExpectQuery("SELECT \\* from pictures").WillReturnRows(rows)

	if got, err := pictureService.GetPictures(); got != nil {
		if err != nil {
			t.Errorf("expected err %+v to be nil", err)
		}
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected result to be %+v, got %+v", expected, got)
		}
	} else {
		t.Errorf("result should not be empty, got err %s", err)
	}
}

func getPictureShouldThrowIfDbThrows(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	expected := errors.New("test error")
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from pictures").WithArgs(id).WillReturnError(expected)

	if _, got := pictureService.GetPicture(id); got != nil {
		if got != expected {
			t.Errorf("error expected %v got %+v", expected, got)
		}
	} else {
		t.Errorf("Error should have been thrown")
	}
}

func getPictureShouldReturnNilIfNoRows(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from pictures").WithArgs(id).WillReturnError(sql.ErrNoRows)

	if got, err := pictureService.GetPicture(id); got == nil {
		if err != nil {
			t.Errorf("expected err %+v to be nil", err)
		}
	} else {
		t.Errorf("result should be nill, got %+v", got)
	}
}

func getPictureShouldReturnPicture(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	rows := buildPictureRows()
	expected := addPictureRows(rows, 1)[0]
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from pictures").WithArgs(id).WillReturnRows(rows)

	if got, err := pictureService.GetPicture(id); got != nil {
		if err != nil {
			t.Errorf("expected err %+v to be nil", err)
		}
		if reflect.DeepEqual(got, expected) {
			t.Errorf("expected result to be %+v, got %+v", expected, got)
		}
	} else {
		t.Errorf("result should not be empty")
	}
}

func createPictureShouldReturnPictureAndCommit(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	picture := createFakePicture()
	user := servertesting.CreateFakeUser()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO pictures (.*)").WithArgs(
		picture.StartDate, picture.Duration, picture.EndDate, picture.Status, user.UserID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if got, err := pictureService.CreatePicture(picture, user); err == nil {
		if got == nil {
			t.Errorf("expected response not to be null")
		}
		if picture != got {
			t.Errorf("expected result to be %+v, got %+v", picture, got)
		}
	} else {
		t.Errorf("error should be null, got %w", err)
	}
}

func createPictureShouldThrowIfPictureIsNull(t *testing.T) {

}

func createPictureShouldThrowIfDbThrows(t *testing.T) {

}

func updatePictureShouldReturnPictureAndCommit(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	picture := createFakePicture()
	id := uint(1)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE pictures SET (.*)").WithArgs(
		picture.StartDate, picture.Duration, picture.EndDate, picture.Status, picture.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if got, err := pictureService.UpdatePicture(id, picture); err == nil {
		if got == nil {
			t.Errorf("expected response not to be null")
		}
		if picture != got {
			t.Errorf("expected result to be %+v, got %+v", picture, got)
		}
	} else {
		t.Errorf("error should be null, got %w", err)
	}
}

func updatePictureShouldThrowIfDbThrows(t *testing.T) {
	dbMock, mock := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	picture := createFakePicture()
	id := uint(1)
	expectedError := errors.New("test error")
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE pictures SET (.*)").WithArgs(
		picture.StartDate, picture.Duration, picture.EndDate, picture.Status, picture.ID,
	).WillReturnError(expectedError)
	mock.ExpectRollback()
	if got, err := pictureService.UpdatePicture(id, picture); err != nil {
		if got != nil {
			t.Errorf("expected response to be null")
		}
		if err != expectedError {
			t.Errorf("expected err to be %+v, got %+v", expectedError, err)
		}
	} else {
		t.Errorf("error should not be null")
	}
}

func updatePictureShouldThrowIfPictureIsNull(t *testing.T) {
	dbMock, _ := mockPicturesDB(t)
	pictureService := NewPictureService(dbMock)
	id := uint(1)
	if got, err := pictureService.UpdatePicture(id, nil); err != nil {
		if got != nil {
			t.Errorf("expected response to be null")
		}
		if err != PictureNullError {
			t.Errorf("expected err to be %+v, got %+v", PictureNullError, err)
		}
	} else {
		t.Errorf("error should not be null")
	}
}

func deletePictureShouldThrowIfDbThrows(t *testing.T) {

}

func deletePictureShouldExecuteAndCommit(t *testing.T) {

}

func mockPicturesDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func buildPictureRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "start_date", "duration", "end_date", "status", "user_id"})
}

func addPictureRows(rows *sqlmock.Rows, numRows uint) []models.Picture {
	var pictures []models.Picture
	var picture *models.Picture
	for ; numRows > 0; numRows-- {
		picture = createFakePicture()

		pictures = append(pictures, *picture)
		rows.AddRow(picture.ID, picture.StartDate, picture.Duration, picture.EndDate, picture.Status, picture.UserId)
	}
	return pictures
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
