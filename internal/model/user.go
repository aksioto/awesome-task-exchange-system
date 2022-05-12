package model

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null"
	"time"
)

const UserKey string = "User"

type User struct {
	Id        int       `json:"id" db:"id"`
	PublicID  uuid.UUID `json:"public_id" db:"public_id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Name      string    `json:"name" db:"name"`
	RoleID    int       `json:"role_id" db:"role_id"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt null.Time `json:"updated_at,omitempty" db:"updated_at"`
}
