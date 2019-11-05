package services

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/models"
)

type TaskService interface {
	GetTasks() ([]models.Task, error)
	GetTask(uint) (*models.Task, error)
}

type taskService struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) TaskService {
	return &taskService{
		DB: db,
	}
}

func (s *taskService) GetTasks() (tasks []models.Task, err error) {
	err = s.DB.Select(&tasks, "SELECT * from tasks")
	if tasks == nil {
		tasks = []models.Task{}
	}
	return
}

func (s *taskService) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	err := s.DB.Get(&task, "SELECT * from tasks where id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, nil
}

func (s *taskService) UpdateTask(id uint) (*models.Task, error) {
	var task models.Task
	err := s.DB.Get(&task, "SELECT * from tasks where id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, nil
}
