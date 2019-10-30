package models

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type Task struct {
	ID          uint16      `json:"id"`
	Title       string      `json:"title"  db:"title"`
	Description null.String `json:"description"  db:"description"`
	Status      string      `json:"status"  db:"status"`
	Priority    int         `json:"priority"  db:"priority"`
	StartDate   null.Time   `json:"start_date" db:"start_date"`
	DueDate     null.Time   `json:"due_date" db:"due_date"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
}
