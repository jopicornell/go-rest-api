package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/internals/errors"
	taskService "github.com/jopicornell/go-rest-api/pkg/api/tasks/services"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"log"
	"net/http"
	"strconv"
)

var tService = &taskService.TaskService{
	DB: database.GetDB(),
}

func GetTasksHandler(w http.ResponseWriter, _ *http.Request) (interface{}, error) {
	tasks, err := tService.GetTasks()
	if err != nil {
		log.Println(fmt.Errorf("error getting tasks: %w", err))
		return nil, errors.InternalServerError
	}
	return tasks, nil
}

func GetOneTaskHandler(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(fmt.Errorf("error converting id: %s to integer", mux.Vars(r)["id"]))
		return nil, errors.InternalServerError
	}
	task, err := tService.GetTask(uint(id))
	if err != nil {
		log.Println(fmt.Errorf("error getting task(%d): %w", id, err))
		return nil, err
	}
	if task == nil {
		return nil, errors.NotFound
	}
	return task, nil
}
