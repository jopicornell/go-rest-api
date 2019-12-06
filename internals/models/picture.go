package models

import (
	"time"
)

const StatusPending = "pending"

type Picture struct {
	ID         uint16        `json:"id" db:"id"`
	StartDate  time.Time     `json:"startDate"  db:"start_date"`
	Duration   time.Duration `json:"duration"  db:"duration"`
	EndDate    time.Time     `json:"endDate"  db:"end_date"`
	Status     string        `json:"status"  db:"status"`
	UserId     int           `json:"userId" db:"user_id"`
	ResourceId int           `json:"resourceId" db:"resource_id"`
}
