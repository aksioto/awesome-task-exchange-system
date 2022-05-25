package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	Username string    `json:"username,omitempty"`
	RoleID   int       `json:"role_id,omitempty"`
	PublicID uuid.UUID `json:"public_id,omitempty"`
	jwt.StandardClaims
}

type ResponseMessage struct {
	Msg    string  `json:"msg,omitempty"`
	Code   int     `json:"code,omitempty"`
	Claims *Claims `json:"claims,omitempty"`
}

const (
	ROLE_ADMIN      = 1
	ROLE_ACCOUNTANT = 2
	ROLE_MANAGER    = 3
	ROLE_USER       = 4
)
