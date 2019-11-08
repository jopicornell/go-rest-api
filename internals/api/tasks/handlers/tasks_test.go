package handlers

import (
	"encoding/json"
	goErrors "errors"
	"github.com/bxcodec/faker/v3"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"net/http/httptest"
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

func (ts *TaskServiceMock) DeleteTask(id uint) error {
	if ts.errorToThrow != nil {
		return ts.errorToThrow
	}
	return nil
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
	t.Run("should throw a NotFound if service says so", notFoundIfServiceFoundsNothing)
	t.Run("should return task returned by the service", returnTaskByService)
}

func TestTaskHandler_GetTasksHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", internalErrorIfServiceFailsReturningTasks)
	t.Run("should return tasks returned by the service", returnTasksByService)
}

func TestTaskHandler_UpdateTaskHandler(t *testing.T) {
	t.Run("should throw a NotFound if service says so", updateTaskShouldThrowNotFound)
	t.Run("should throw InternalServer if some error is raised by the service", updateTaskShouldThrowIfServiceFails)
	t.Run("should return task updated by the service", updateTaskShouldReturnUpdatedTask)
}

func TestTaskHandler_CreateTaskHandler(t *testing.T) {
	t.Run("should throw InternalServer if some error is raised by the service", createTaskShouldThrowIfServiceFails)
	t.Run("should return task created by the service", createTaskShouldReturnCreatedTask)
}

func panicErrorServiceMissing(t *testing.T) {
	mock := &TaskHandlerMock{}
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/tasks/1", nil),
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
	mock.GetOneTaskHandler(context)
}

func returnTaskByService(t *testing.T) {
	mock := &TaskHandlerMock{}
	expected := &models.Task{
		ID:          0,
		Title:       "",
		Description: null.String{},
		Date:        null.Time{},
	}
	mock.taskService = &TaskServiceMock{task: expected}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/tasks/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetOneTaskHandler(context)
	var got *models.Task
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed: %w %s", err, recorder.Body.String())
	}
	if *got != *expected {
		t.Errorf("expected %+v got %+v", expected, got)
	}
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d got %d", recorder.Code, http.StatusOK)
	}
}

func returnTasksByService(t *testing.T) {
	mock := &TaskHandlerMock{}
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
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/tasks/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetTasksHandler(context)
	var got []models.Task
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed (%s): %w", recorder.Body.String(), err)
	}
	if len(got) != len(expected) {
		t.Errorf("expected length of result (%d) to be %d", len(got), len(expected))
	}

}

func internalErrorIfServiceFailsReturningTask(t *testing.T) {
	mock := &TaskHandlerMock{}
	var expected *models.Task = nil
	mock.taskService = &TaskServiceMock{
		task:         expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/tasks/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetOneTaskHandler(context)
	if recorder.Code != 500 {
		t.Errorf("expected status code to be 500 got %+v", recorder.Code)
	}
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v (%s)", recorder.Body.Len(), recorder.Body.String())
	}
}

func internalErrorIfServiceFailsReturningTasks(t *testing.T) {
	mock := &TaskHandlerMock{}
	var expected []models.Task = nil
	mock.taskService = &TaskServiceMock{
		tasks:        expected,
		errorToThrow: goErrors.New("test error"),
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/tasks/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetTasksHandler(context)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status code to be %d got %d", http.StatusInternalServerError, recorder.Code)
	}
}

func notFoundIfServiceFoundsNothing(t *testing.T) {
	mock := &TaskHandlerMock{}
	mock.taskService = &TaskServiceMock{}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("GET", "/tasks/1", nil),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.GetOneTaskHandler(context)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %+v got %+v", http.StatusNotFound, recorder.Code)
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
	expected := createFakeTask()
	mock.taskService = &TaskServiceMock{
		task:         expected,
		errorToThrow: nil,
	}
	var taskJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(taskJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.UpdateTaskHandler(context)
	var got *models.Task
	if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
		t.Errorf("unmarshalling failed: %w", err)
	}
	if *got != *expected {
		t.Errorf("expected %+v got %+v", expected, got)
	}
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d got %d", recorder.Code, http.StatusOK)
	}

}

func updateTaskShouldThrowIfServiceFails(t *testing.T) {
	mock := &TaskHandlerMock{}
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
	context := servertesting.NewContext(
		httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(taskJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.UpdateTaskHandler(context)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status code to be %d got %d", http.StatusInternalServerError, recorder.Code)
	}

}

func updateTaskShouldThrowNotFound(t *testing.T) {
	mock := &TaskHandlerMock{}
	task := createFakeTask()
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
	context := servertesting.NewContext(
		httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(taskJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.UpdateTaskHandler(context)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %+v got %+v", http.StatusNotFound, recorder.Code)
	}
}

func createTaskShouldReturnCreatedTask(t *testing.T) {
	mock := &TaskHandlerMock{}
	expected := createFakeTask()
	mock.taskService = &TaskServiceMock{
		task:         expected,
		errorToThrow: nil,
	}
	var taskJSON string
	if marshallResult, err := json.Marshal(expected); err == nil {
		taskJSON = string(marshallResult)
	} else {
		t.Errorf("error marshalling task %w", err)
	}
	recorder := httptest.NewRecorder()
	context := servertesting.NewContext(
		httptest.NewRequest("POST", "/tasks", strings.NewReader(taskJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.CreateTaskHandler(context)
	var got *models.Task
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

func createTaskShouldThrowIfServiceFails(t *testing.T) {
	mock := &TaskHandlerMock{}

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
	context := servertesting.NewContext(
		httptest.NewRequest("POST", "/tasks", strings.NewReader(taskJSON)),
		recorder,
		map[string]string{
			"id": "1",
		},
	)
	mock.CreateTaskHandler(context)
	if recorder.Body.Len() != 0 {
		t.Errorf("expected length to be 0 got %+v", recorder.Body.Len())
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status code to be %d got %d", http.StatusInternalServerError, recorder.Code)
	}
}

func createFakeTask() *models.Task {
	return &models.Task{
		ID:          uint16(faker.UnixTime()),
		Title:       faker.Sentence(),
		Description: null.NewString(faker.Sentence(), true),
		Completed:   false,
		Date:        null.NewTime(time.Now().Round(time.Nanosecond), true),
	}
}
