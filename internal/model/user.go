package model

import (
	"github.com/google/uuid"
	"time"
)

//easyjson:json
type User struct {
	ID           uuid.UUID `json:"-"`
	Role         Role      `json:"role"`
	Name         string    `json:"name"`
	Login        string    `json:"login"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

//easyjson:json
type UpdateUserPayload struct {
	ID   uuid.UUID `json:"-"`
	Name string    `json:"name" validate:"required,min=2,max=30"`
}
