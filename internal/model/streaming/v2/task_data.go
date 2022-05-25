package v2

import "github.com/google/uuid"

type TaskData struct {
	PublicTaskID uuid.UUID `json:"public_task_id,omitempty"`
	Title        string    `json:"title,omitempty"`
	JiraID       string    `json:"jira_id,omitempty"`
	Description  string    `json:"description,omitempty"`
}
