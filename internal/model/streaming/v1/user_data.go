package v1

import "github.com/google/uuid"

type UserData struct {
	PublicID uuid.UUID `json:"public_id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	RoleID   int       `json:"role_id"`
}
