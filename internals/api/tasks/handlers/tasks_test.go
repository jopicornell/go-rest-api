package handlers

import (
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"gopkg.in/guregu/null.v3"
	"net/http/httptest"
	"testing"
)

type TaskHandlerMock struct {
	TaskHandler
}

type TaskServiceMock struct {
	task  *models.Task
	tasks []models.Task
}

func (ts *TaskServiceMock) GetTask(id uint) (*models.Task, error) {
	return ts.task, nil
}

func (ts *TaskServiceMock) GetTasks() ([]models.Task, error) {
	return ts.tasks, nil
}

func TestTaskHandler_GetOneTaskHandler(t *testing.T) {
	t.Run("should throw error if service is missing", panicErrorServiceMissing)
	t.Run("should throw an InternalServerError if id is invalid", failIfTaskIdIsInvalid)
	t.Run("should throw a NotFound if service says so", notFoundIfServiceFoundsNothing)
	t.Run("should return task returned by the service", returnTaskByService)
}

func TestTaskHandler_GetTasksHandler(t *testing.T) {
	t.Run("should throw error if service is missing", panicErrorServiceMissing)
	t.Run("should throw an InternalServerError if id is invalid", failIfTaskIdIsInvalid)
	t.Run("should throw a NotFound if service says so", notFoundIfServiceFoundsNothing)
	t.Run("should return task returned by the service", returnTaskByService)
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
