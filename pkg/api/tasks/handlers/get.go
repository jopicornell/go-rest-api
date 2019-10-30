package handlers

import (
	"github.com/gorilla/mux"
	taskService "github.com/jopicornell/go-rest-api/pkg/api/tasks/services"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"github.com/jopicornell/go-rest-api/pkg/util"
	"net/http"
	"strconv"
)

var tService = &taskService.Service{
	DB: database.GetDB(),
}

func GetTasksHandler(w http.ResponseWriter, _ *http.Request) {
	tasks, err := tService.GetTasks()
	if err != nil {
		util.WriteInternalErrorToResponse(w, err)
		return
	}
	util.WriteToJson(w, tasks)
}

func GetOneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		util.WriteInternalErrorToResponse(w, err)
	}
	task, err := tService.GetTask(uint(id))
	if err != nil {
		util.WriteInternalErrorToResponse(w, err)
		return
	}
	if task == nil {
		w.WriteHeader(404)
		return
	}
	util.WriteToJson(w, task)
}
