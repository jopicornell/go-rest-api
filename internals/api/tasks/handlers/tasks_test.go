package handlers

import (
	"encoding/json"
	goErrors "errors"
	"github.com/bxcodec/faker/v3"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/models"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"gopkg.in/guregu/null.v3"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type TaskHandlerMock struct {
	TaskHandler
}

type TaskServiceMock struct {
	task         *models.Task
	tasks        []models.Task
	errorToThrow error
}

func (ts *TaskServiceMock) UpdateTask(id uint, task *models.Task) (*models.Task, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.task, nil
}

func (ts *TaskServiceMock) CreateTask(task *models.Task) (*models.Task, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.task, nil
}

func (ts *TaskServiceMock) GetTask(id uint) (*models.Task, error) {
	if ts.errorToThrow != nil {
		return nil, ts.errorToThrow
	}
	return ts.task, nil
}

func (ts *TaskServiceMock) GetTasks() ([]models.Task, error) {
	if ts.errorToThrow != nil {
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

func TestTaskHandler_UpdateTaskHandler(t *testing.T) {
	t.Run("should throw error if id is invalid", updateTaskShouldFailIfTaskIdIsInvalid)
	t.Run("should throw error if body is invalid", updateTaskShouldFailIfBodyIdIsInvalid)
	t.Run("should throw error if json is invalid", updateTaskShouldFailIfJsonIdIsInvalid)
	t.Run("should throw a NotFound if service says so", updateTaskShouldThrowNotFound)
	t.Run("should throw InternalServer if some error is raised by the service", updateTaskShouldThrowIfServiceFails)
	t.Run("should return task updated by the service", updateTaskShouldReturnUpdatedTask)
}

func TestTaskHandler_CreateTaskHandler(t *testing.T) {
	t.Run("should throw error if body is invalid", createTaskShouldFailIfBodyIdIsInvalid)
	t.Run("should throw error if json is invalid", createTaskShouldFailIfJsonIdIsInvalid)
	t.Run("should throw InternalServer if some error is raised by the service", createTaskShouldThrowIfServiceFails)
	t.Run("should return task created by the service", createTaskShouldReturnCreatedTask)
}

func panicErrorServiceMissing(t *testing.T) {
	mock := &TaskHandlerMock{}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
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
	request := servertesting.NewRequest(httptest.NewRequest("GET", "/tasks/asd", nil), map[string]string{
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
	mock.Server = servertesting.Initialize(&config.Config{})
	expected := &models.Task{
		ID:          0,
		Title:       "",
		Description: null.String{},
		Date:        null.Time{},
	}
	mock.taskService = &TaskServiceMock{task: expected}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
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
	mock.Server = servertesting.Initialize(&config.Config{})
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
	request := servertesting.NewRequest(httptest.NewRequest("GET", "/tasks", nil), map[string]string{})
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
	mock.Server = servertesting.Initialize(&config.Config{})
	var expected *models.Task = nil
	mock.taskService = &TaskServiceMock{
		task:         expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
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
	mock.Server = servertesting.Initialize(&config.Config{})
	var expected []models.Task = nil
	mock.taskService = &TaskServiceMock{
		tasks:        expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(httptest.NewRequest("GET", "/tasks", nil), map[string]string{})
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
	mock.Server = servertesting.Initialize(&config.Config{})
	mock.taskService = &TaskServiceMock{}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(httptest.NewRequest("GET", "/tasks/1", nil), map[string]string{
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
	serverMock := &servertesting.ServerMock{
		Config: config.Config{},
	}
	taskHandler := New(serverMock)
	if taskHandler == nil {
		t.Errorf("task handler should not be null")
		return
	}
	if taskHandler.taskService == nil {
		t.Errorf("task handler created without the service")
	}
}

func updateTaskShouldReturnUpdatedTask(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = servertesting.Initialize(&config.Config{})
	task := createFakeTask()
	mock.taskService = &TaskServiceMock{
		task:         task,
		errorToThrow: nil,
	}
	var taskJSON string
	if marshallResult, err := json.Marshal(task); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/tasks/1", strings.NewReader(taskJSON)),
		map[string]string{
			"id": "1",
		},
	)
	if got, err := mock.UpdateTaskHandler(recorder, request); err == nil {
		if got == nil {
			t.Errorf("expected result not to be nil")
		}
		if got != task {
			t.Errorf("expected result to be %+v, got %+v", task, got)
		}
	} else {
		t.Errorf("expected %s to be nil", err)
	}
}

func updateTaskShouldThrowIfServiceFails(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = servertesting.Initialize(&config.Config{})
	task := createFakeTask()
	errorToThrow := goErrors.New("test error")
	mock.taskService = &TaskServiceMock{
		task:         task,
		errorToThrow: errorToThrow,
	}
	var taskJSON string
	if marshallResult, err := json.Marshal(task); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/tasks/1", strings.NewReader(taskJSON)),
		map[string]string{
			"id": "1",
		},
	)
	if got, err := mock.UpdateTaskHandler(recorder, request); err != nil {
		if got != nil {
			t.Errorf("expected result to be nil")
		}
		if err != errorToThrow {
			t.Errorf("expected result to be %+v, got %+v", task, got)
		}
	} else {
		t.Errorf("expected err not to be nil")
	}
}

func updateTaskShouldFailIfTaskIdIsInvalid(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.taskService = &TaskServiceMock{}
	recorder := httptest.NewRecorder()
	task := createFakeTask()
	var taskJSON string
	if marshallResult, err := json.Marshal(task); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}

	request := servertesting.NewRequest(httptest.NewRequest("PUT", "/tasks/asd", strings.NewReader(taskJSON)), map[string]string{
		"id": "asd",
	})
	res, err := mock.UpdateTaskHandler(recorder, request)
	if res != nil {
		t.Errorf("expected response to be nil, got %s", res)
	}
	if err != errors.InternalServerError {
		t.Error("expected err to be an internal server error")
	}
}

func updateTaskShouldFailIfJsonIdIsInvalid(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = servertesting.Initialize(&config.Config{})
	taskJSON := "Invalid json;{{}"

	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/tasks/1", strings.NewReader(taskJSON)),
		map[string]string{
			"id": "1",
		},
	)
	if got, err := mock.UpdateTaskHandler(recorder, request); err != nil {
		if got != nil {
			t.Errorf("expected result to be nil")
		}
		if err != errors.BadRequest {
			t.Errorf("expected error to be %+v, got %+v", errors.BadRequest, err)
		}
	} else {
		t.Errorf("expected err not to be nil")
	}
}

func updateTaskShouldFailIfBodyIdIsInvalid(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = servertesting.Initialize(&config.Config{})
	task := createFakeTask()
	var taskJSON string
	if marshallResult, err := json.Marshal(task); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/tasks/1", strings.NewReader(taskJSON)),
		map[string]string{
			"id": "1",
		},
	)
	request.ThrowError = goErrors.New("error parsing body")
	if got, err := mock.UpdateTaskHandler(recorder, request); err != nil {
		if got != nil {
			t.Errorf("expected result to be nil")
		}
		if err != errors.BadRequest {
			t.Errorf("expected error to be %+v, got %+v", errors.BadRequest, err)
		}
	} else {
		t.Errorf("expected err not to be nil")
	}
}
func updateTaskShouldThrowNotFound(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = servertesting.Initialize(&config.Config{})
	task := createFakeTask()
	errorToThrow := errors.NotFound
	mock.taskService = &TaskServiceMock{
		task:         nil,
		errorToThrow: nil,
	}
	var taskJSON string
	if marshallResult, err := json.Marshal(task); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/tasks/1", strings.NewReader(taskJSON)),
		map[string]string{
			"id": "1",
		},
	)
	if got, err := mock.UpdateTaskHandler(recorder, request); err != nil {
		if got != nil {
			t.Errorf("expected result to be nil")
		}
		if err != errorToThrow {
			t.Errorf("expected result to be %+v, got %+v", task, got)
		}
	} else {
		t.Errorf("expected err not to be nil")
	}
}

func createTaskShouldReturnCreatedTask(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = servertesting.Initialize(&config.Config{})
	task := createFakeTask()
	mock.taskService = &TaskServiceMock{
		task:         task,
		errorToThrow: nil,
	}
	var taskJSON string
	if marshallResult, err := json.Marshal(task); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(
		httptest.NewRequest("GET", "/tasks/1", strings.NewReader(taskJSON)), map[string]string{},
	)
	if got, err := mock.CreateTaskHandler(recorder, request); err == nil {
		if got == nil {
			t.Errorf("expected result not to be null")
		}
	} else {
		t.Errorf("expected err to be nil, got %w", err)
	}
}

func createTaskShouldThrowIfServiceFails(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.Server = servertesting.Initialize(&config.Config{})

	task := createFakeTask()
	var taskJSON string
	if marshallResult, err := json.Marshal(task); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	mock.taskService = &TaskServiceMock{
		task:         nil,
		tasks:        nil,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	request := servertesting.NewRequest(
		httptest.NewRequest("POST", "/tasks", strings.NewReader(taskJSON)), map[string]string{},
	)
	if got, err := mock.CreateTaskHandler(recorder, request); err != nil {
		if got != nil {
			t.Errorf("expected result to be nil")
		}
		if err != errors.InternalServerError {
			t.Errorf("expected result to be %+v, got %+v", errors.InternalServerError, err)
		}
	} else {
		t.Errorf("expected err not to be nil")
	}
}

func createTaskShouldFailIfJsonIdIsInvalid(t *testing.T) {

}

func createTaskShouldFailIfBodyIdIsInvalid(t *testing.T) {

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
