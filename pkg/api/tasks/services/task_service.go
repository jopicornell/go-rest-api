package services

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"github.com/jopicornell/go-rest-api/pkg/models"
)

type TaskService struct {
	DB *sqlx.DB
}

func (s *TaskService) GetTasks() (tasks []models.Task, err error) {
	err = database.GetDB().Select(&tasks, "SELECT * from tasks")
	if tasks == nil {
		tasks = []models.Task{}
	}
	return
}

func (s *TaskService) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	err := database.GetDB().Get(&task, "SELECT * from tasks where id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, nil
}

func (s *TaskService) UpdateTask(id uint) (*models.Task, error) {
	var task models.Task
	err := database.GetDB().Get(&task, "SELECT * from tasks where id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, nil
}
