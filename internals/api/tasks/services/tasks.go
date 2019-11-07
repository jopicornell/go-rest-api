package services

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/models"
)

type TaskService interface {
	GetTasks() ([]models.Task, error)
	GetTask(uint) (*models.Task, error)
	UpdateTask(uint, *models.Task) (*models.Task, error)
	CreateTask(*models.Task) (*models.Task, error)
}

type taskService struct {
	db *sqlx.DB
}

var TaskNullError = errors.New("task should not be null")

func New(db *sqlx.DB) TaskService {
	return &taskService{
		db: db,
	}
}

func (s *taskService) GetTasks() (tasks []models.Task, err error) {
	tasks = []models.Task{}
	if err = s.db.Select(&tasks, "SELECT * from tasks"); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *taskService) GetTask(id uint) (*models.Task, error) {
	var task models.Task
	if err := s.db.Get(&task, "SELECT * from tasks where id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (s *taskService) CreateTask(task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, TaskNullError
	}
	insertQuery := "INSERT INTO tasks (title, description, date, completed) VALUES (?, ?, ?, 0)"
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(insertQuery, task.Title, task.Description, task.Date); err == nil {
		err = tx.Commit()
		return task, err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, err
		}
		return nil, err
	}
}

func (s *taskService) UpdateTask(id uint, task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, TaskNullError
	}
	updateQuery := "UPDATE tasks SET title=?, description=?, completed=? where task_id = ?"
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(updateQuery, task.Title, task.Description, task.Completed, task.ID); err == nil {
		err = tx.Commit()
		return task, err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, err
		}
		return nil, err
	}
}

func (s *taskService) DeleteTask(id uint) (err error) {
	deleteQuery := "DELETE FROM tasks WHERE id = ?"
	tx, err := s.db.Beginx()
	if _, err := tx.Exec(deleteQuery); err == nil {
		err = tx.Commit()
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}
}
