package handlers

import (
	goErrors "errors"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"gopkg.in/guregu/null.v3"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TaskHandlerMock struct {
	TaskHandler
}

type TaskServiceMock struct {
	task         *models.Task
	tasks        []models.Task
	throwError   bool
	errorToThrow error
}

func (ts *TaskServiceMock) GetTask(id uint) (*models.Task, error) {
	if ts.throwError {
		return nil, ts.errorToThrow
	}
	return ts.task, nil
}

func (ts *TaskServiceMock) GetTasks() ([]models.Task, error) {
	if ts.throwError {
		return nil, ts.errorToThrow
	}
	return ts.tasks, nil
}

func TestTaskHandler_New(t *testing.T) {
	t.Run("should construct a new TaskHandler given the server", shouldReturnConstructedHandler)
}

func TestTaskHandler_GetOneTaskHandler(t *testing.T) {
	t.Run("should throw error if service is missing", panicErrorServiceMissing)
	t.Run("should throw InternalServerError if service fails", internalErrorIfServiceFailsReturningTask)
	t.Run("should throw an InternalServerError if id is invalid", failIfTaskIdIsInvalid)
	t.Run("should throw a NotFound if service says so", notFoundIfServiceFoundsNothing)
	t.Run("should return task returned by the service", returnTaskByService)
}

func TestTaskHandler_GetTasksHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", internalErrorIfServiceFailsReturningTasks)
	t.Run("should return tasks returned by the service", returnTasksByService)
}

func panicErrorServiceMissing(t *testing.T) {
	mock := &TaskHandlerMock{}
	recorder := httptest.NewRecorder()
	request := servertesting.New(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
		"id": "1",
	})
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("function did not panic")
		}
	}()

	_, _ = mock.GetOneTaskHandler(recorder, request)
}

func failIfTaskIdIsInvalid(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.taskService = &TaskServiceMock{}
	recorder := httptest.NewRecorder()
	request := servertesting.New(httptest.NewRequest("GET", "/tasks/asd", nil), map[string]string{
		"id": "asd",
	})
	res, err := mock.GetOneTaskHandler(recorder, request)
	if res != nil {
		t.Errorf("expected response to be nil, got %s", res)
	}
	if err != errors.InternalServerError {
		t.Error("expected err to be an internal server error")
	}
}

func returnTaskByService(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = &server.Server{
		Config: config.Config{},
	}
	expected := &models.Task{
		ID:          0,
		Title:       "",
		Description: null.String{},
		Date:        null.Time{},
	}
	mock.taskService = &TaskServiceMock{task: expected}
	recorder := httptest.NewRecorder()
	request := servertesting.New(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
		"id": "1",
	})
	if res, err := mock.GetOneTaskHandler(recorder, request); err == nil {
		if res != expected {
			t.Errorf("expected %+v got %+v", expected, res)
		}
	} else {
		t.Errorf("expected %s to be nil", err)
	}
}

func returnTasksByService(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = &server.Server{
		Config: config.Config{},
	}
	expected := []models.Task{
		{
			ID:          0,
			Title:       "",
			Description: null.String{},
			Date:        null.Time{},
		},
		{
			ID:          0,
			Title:       "",
			Description: null.String{},
			Date:        null.Time{},
		},
	}
	mock.taskService = &TaskServiceMock{tasks: expected}
	recorder := httptest.NewRecorder()
	request := servertesting.New(httptest.NewRequest("GET", "/tasks", nil), map[string]string{})
	if res, err := mock.GetTasksHandler(recorder, request); err == nil {
		resValue := reflect.ValueOf(res)
		if resValue.Kind() != reflect.Slice {
			t.Errorf("expected %+v to be a slice", res)
		}
		if resValue.Len() != len(expected) {
			t.Errorf("expected length of result (%d) to be %d", resValue.Len(), len(expected))
		}
	} else {
		t.Errorf("expected %s to be nil", err)
	}
}

func internalErrorIfServiceFailsReturningTask(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = &server.Server{
		Config: config.Config{},
	}
	var expected *models.Task = nil
	mock.taskService = &TaskServiceMock{
		task:         expected,
		throwError:   true,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	request := servertesting.New(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
		"id": "1",
	})
	if res, err := mock.GetOneTaskHandler(recorder, request); err != nil {
		if res != nil {
			t.Errorf("expected %+v got %+v", expected, res)
		}
	} else {
		t.Errorf("expected %s to not be nil", err)
	}
}

func internalErrorIfServiceFailsReturningTasks(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = &server.Server{
		Config: config.Config{},
	}
	var expected []models.Task = nil
	mock.taskService = &TaskServiceMock{
		tasks:        expected,
		throwError:   true,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	request := servertesting.New(httptest.NewRequest("GET", "/tasks", nil), map[string]string{})
	if res, err := mock.GetTasksHandler(recorder, request); err != nil {
		if res != nil {
			t.Errorf("expected nil got %+v", res)
		}
	} else {
		t.Errorf("expected %s to not be nil", err)
	}
}

func notFoundIfServiceFoundsNothing(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = &server.Server{
		Config: config.Config{},
	}
	mock.taskService = &TaskServiceMock{}
	recorder := httptest.NewRecorder()
	request := servertesting.New(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
		"id": "1",
	})
	expected := errors.NotFound
	if res, err := mock.GetOneTaskHandler(recorder, request); err != nil {
		if err != expected {
			t.Errorf("expected %+v got %+v", expected, res)
		}
	} else {
		t.Errorf("expected %s to be nil", res)
	}
}

func shouldReturnConstructedHandler(t *testing.T) {
	serverMock := server.Initialize()
	taskHandler := New(serverMock)
	if taskHandler == nil {
		t.Errorf("task handler should not be null")
		return
	}
	if taskHandler.taskService == nil {
		t.Errorf("task handler created without the service")
	}
}
