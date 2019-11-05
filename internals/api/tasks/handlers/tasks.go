package handlers

import (
	"fmt"
	"github.com/jopicornell/go-rest-api/internals/api/tasks/services"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"log"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	server.Handler
	taskService services.TaskService
}

func New(s *server.Server) *TaskHandler {
	return &TaskHandler{
		taskService: services.New(s.GetRelationalDatabase()),
	}
}

func (s *TaskHandler) GetTasksHandler(w http.ResponseWriter, _ server.Request) (interface{}, error) {
	tasks, err := s.taskService.GetTasks()
	if err != nil {
		log.Println(fmt.Errorf("error getting tasks: %w", err))
		return nil, errors.InternalServerError
	}
	return tasks, nil
}

func (s *TaskHandler) GetOneTaskHandler(_ http.ResponseWriter, r server.Request) (interface{}, error) {
	id, err := strconv.Atoi(r.GetPathParameters()["id"])
	if err != nil {
		log.Println(fmt.Errorf("error converting id: %s to integer", r.GetPathParameters()["id"]))
		return nil, errors.InternalServerError
	}
	task, err := s.taskService.GetTask(uint(id))
	if err != nil {
		log.Println(fmt.Errorf("error getting task(%d): %w", id, err))
		return nil, err
	}
	if task == nil {
		return nil, errors.NotFound
	}
	return task, nil
}
