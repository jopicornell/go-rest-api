package services

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"github.com/jopicornell/go-rest-api/pkg/models"
)

type Service struct {
	DB *sqlx.DB
}

func (s *Service) GetTasks() (tasks []models.Task, err error) {
	err = database.GetDB().Select(&tasks, "SELECT * from tasks")
	if tasks == nil {
		tasks = []models.Task{}
	}
	return
}

func (s *Service) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	err := database.GetDB().Get(&task, "SELECT * from tasks where id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, nil
}
