package services

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/models"
	"gopkg.in/guregu/null.v3"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	dbMock, _ := mockDB(t)
	taskService := New(dbMock)
	if taskService == nil {
		t.Errorf("New service should not be null")
	}
	reflectServiceDb := reflect.ValueOf(taskService).Elem().FieldByName("db")
	if reflectServiceDb.IsNil() {
		t.Errorf("db field in service should not be nil")
	}
}

func TestTaskService_GetTasks(t *testing.T) {
	t.Run("should throw if db throws", getTasksShouldThrowIfDbThrows)
	t.Run("should return empty slice if no rows", getTasksShouldReturnEmptySliceIfNoRows)
	t.Run("should return list of tasks if all went ok", getTasksShouldReturnListOfTasks)
}

func TestTaskService_GetTask(t *testing.T) {
	t.Run("should throw if db throws", getTaskShouldThrowIfDbThrows)
	t.Run("should return nil if no rows", getTaskShouldReturnNilIfNoRows)
	t.Run("should return task if all went ok", getTaskShouldReturnTask)
}

func TestTaskService_UpdateTask(t *testing.T) {
	t.Run("should throw if db throws and rollback", updateTaskShouldThrowIfDbThrows)
	t.Run("should throw if task to updateis null", updateTaskShouldThrowIfTaskIsNull)
	t.Run("should return updated task and commit", updateTaskShouldReturnTaskAndCommit)
}

func TestTaskService_DeleteTask(t *testing.T) {
	t.Run("should throw if db throws and rollback", deleteTaskShouldThrowIfDbThrows)
	t.Run("should return no error and commit", deleteTaskShouldExecuteAndCommit)
}

func getTasksShouldThrowIfDbThrows(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	expected := errors.New("test error")
	mock.ExpectQuery("SELECT \\* from tasks").WillReturnError(expected)

	if _, got := taskService.GetTasks(); got != nil {
		if got != expected {
			t.Errorf("error expected %v got %+v", expected, got)
		}
	} else {
		t.Errorf("Error should have been thrown")
	}
}

func getTasksShouldReturnEmptySliceIfNoRows(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	mock.ExpectQuery("SELECT \\* from tasks").WillReturnRows(&sqlmock.Rows{})

	if got, err := taskService.GetTasks(); got != nil {
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

func getTasksShouldReturnListOfTasks(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	rows := buildTaskRows()
	expected := addTaskRows(rows, 5)
	mock.ExpectQuery("SELECT \\* from tasks").WillReturnRows(rows)

	if got, err := taskService.GetTasks(); got != nil {
		if err != nil {
			t.Errorf("expected err %+v to be nil", err)
		}
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected result to be %+v, got %+v", expected, got)
		}
	} else {
		t.Errorf("result should not be empty")
	}
}

func getTaskShouldThrowIfDbThrows(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	expected := errors.New("test error")
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from tasks").WithArgs(id).WillReturnError(expected)

	if _, got := taskService.GetTask(id); got != nil {
		if got != expected {
			t.Errorf("error expected %v got %+v", expected, got)
		}
	} else {
		t.Errorf("Error should have been thrown")
	}
}

func getTaskShouldReturnNilIfNoRows(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from tasks").WithArgs(id).WillReturnError(sql.ErrNoRows)

	if got, err := taskService.GetTask(id); got == nil {
		if err != nil {
			t.Errorf("expected err %+v to be nil", err)
		}
	} else {
		t.Errorf("result should be nill, got %+v", got)
	}
}

func getTaskShouldReturnTask(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	rows := buildTaskRows()
	expected := addTaskRows(rows, 1)[0]
	id := uint(1)
	mock.ExpectQuery("SELECT \\* from tasks").WithArgs(id).WillReturnRows(rows)

	if got, err := taskService.GetTask(id); got != nil {
		if err != nil {
			t.Errorf("expected err %+v to be nil", err)
		}
		if expected != *got {
			t.Errorf("expected result to be %+v, got %+v", expected, got)
		}
	} else {
		t.Errorf("result should not be empty")
	}
}

func updateTaskShouldReturnTaskAndCommit(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	task := createFakeTask()
	id := uint(1)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks SET (.*)").WithArgs(
		task.Title, task.Description.ValueOrZero(), task.Completed, task.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if got, err := taskService.UpdateTask(id, task); err == nil {
		if got == nil {
			t.Errorf("expected response not to be null")
		}
		if task != got {
			t.Errorf("expected result to be %+v, got %+v", task, got)
		}
	} else {
		t.Errorf("error should be null, got %w", err)
	}
}

func updateTaskShouldThrowIfDbThrows(t *testing.T) {
	dbMock, mock := mockDB(t)
	taskService := New(dbMock)
	task := createFakeTask()
	id := uint(1)
	expectedError := errors.New("test error")
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks SET (.*)").WithArgs(
		task.Title, task.Description, task.Completed, task.ID,
	).WillReturnError(expectedError)
	mock.ExpectRollback()
	if got, err := taskService.UpdateTask(id, task); err != nil {
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

func updateTaskShouldThrowIfTaskIsNull(t *testing.T) {
	dbMock, _ := mockDB(t)
	taskService := New(dbMock)
	id := uint(1)
	if got, err := taskService.UpdateTask(id, nil); err != nil {
		if got != nil {
			t.Errorf("expected response to be null")
		}
		if err != TaskNullError {
			t.Errorf("expected err to be %+v, got %+v", TaskNullError, err)
		}
	} else {
		t.Errorf("error should not be null")
	}
}

func deleteTaskShouldThrowIfDbThrows(t *testing.T) {

}

func deleteTaskShouldExecuteAndCommit(t *testing.T) {

}

func mockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func buildTaskRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"task_id", "title", "description", "completed", "date"})
}

func addTaskRows(rows *sqlmock.Rows, numRows uint) []models.Task {
	var tasks []models.Task
	var task *models.Task
	for ; numRows > 0; numRows-- {
		task = createFakeTask()

		tasks = append(tasks, *task)
		rows.AddRow(task.ID, task.Title, task.Description, task.Completed, task.Date)
	}
	return tasks
}

func createFakeTask() *models.Task {
	return &models.Task{
		ID:          uint16(faker.UnixTime()),
		Title:       faker.Sentence(),
		Description: null.NewString(faker.Sentence(), true),
		Completed:   false,
		Date:        null.NewTime(time.Now(), true),
	}
}
