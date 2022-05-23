package model

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null"
	"time"
)

type Task struct {
	Id          int         `json:"id" db:"id"`
	TaskID      uuid.UUID   `json:"task_id" db:"task_id"`
	UserID      uuid.UUID   `json:"user_id" db:"user_id"`
	Title       string      `json:"title" db:"title"`
	JiraID      null.String `json:"jira_id" db:"jira_id"`
	Description string      `json:"description" db:"description"`
	Status      int8        `json:"status" db:"status"`
	CreatedAt   time.Time   `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   null.Time   `json:"updated_at,omitempty" db:"updated_at"`
}
