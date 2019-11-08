package handlers

import (
	"fmt"
	"github.com/jopicornell/go-rest-api/internals/api/tasks/services"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"log"
	"net/http"
)

type TaskHandler struct {
	server.Handler
	taskService services.TaskService
}

func New(s server.Server) *TaskHandler {
	return &TaskHandler{
		taskService: services.New(s.GetRelationalDatabase()),
	}
}

func (s *TaskHandler) GetTasksHandler(context server.Context) {
	tasks, err := s.taskService.GetTasks()
	if err != nil {
		log.Println(fmt.Errorf("error getting tasks: %w", err))
		context.Respond(http.StatusInternalServerError)
		return
	}
	context.RespondJSON(http.StatusOK, tasks)
}

func (s *TaskHandler) GetOneTaskHandler(context server.Context) {
	id := context.GetParamUInt("id")
	task, err := s.taskService.GetTask(uint(id))
	if err != nil {
		log.Println(fmt.Errorf("error getting task(%d): %w", id, err))
		context.Respond(http.StatusInternalServerError)
		return
	}
	if task == nil {
		context.Respond(http.StatusNotFound)
		return
	}
	context.RespondJSON(http.StatusOK, task)
}

func (s *TaskHandler) UpdateTaskHandler(context server.Context) {
	id := context.GetParamUInt("id")
	var task *models.Task
	context.GetBodyMarshalled(&task)
	task, err := s.taskService.UpdateTask(uint(id), task)
	if err != nil {
		log.Println(fmt.Errorf("error getting task(%d): %w", id, err))
		context.Respond(http.StatusInternalServerError)
		return
	}
	if task == nil {
		context.Respond(http.StatusNotFound)
		return
	}
	context.RespondJSON(http.StatusOK, task)
}

func (s *TaskHandler) CreateTaskHandler(context server.Context) {
	var task *models.Task
	context.GetBodyMarshalled(&task)
	if task, err := s.taskService.CreateTask(task); err == nil {
		context.RespondJSON(http.StatusCreated, task)
	} else {
		log.Println(fmt.Errorf("error creating task %+v: %w", task, err))
		context.Respond(http.StatusInternalServerError)
	}
}

func (s *TaskHandler) DeleteTaskHandler(context server.Context) {
	id := context.GetParamUInt("id")
	if err := s.taskService.DeleteTask(uint(id)); err == nil {
		context.Respond(http.StatusOK)
	} else {
		log.Println(fmt.Errorf("error deleting task %d: %w", id, err))
		context.Respond(http.StatusInternalServerError)
	}
}
