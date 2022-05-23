package v1

import "github.com/google/uuid"

type TaskData struct {
	PublicTaskID uuid.UUID `json:"public_task_id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
}

type TaskAssigneeData struct {
	PublicUserID uuid.UUID `json:"public_user_id,omitempty"`
	PublicTaskID uuid.UUID `json:"public_task_id,omitempty"`
}
