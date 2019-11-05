package models

import (
	"gopkg.in/guregu/null.v3"
)

type Task struct {
	ID          uint16      `json:"task_id" db:"task_id"`
	Title       string      `json:"title"  db:"title"`
	Description null.String `json:"description"  db:"description"`
	Completed   bool        `json:"completed"  db:"completed"`
	Date        null.Time   `json:"date" db:"date"`
}
