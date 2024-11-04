package model

import (
	"github.com/google/uuid"
)

//easyjson:json
type SignInPayload struct {
	Login    string `json:"login" validate:"required,min=6,max=30"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

//easyjson:json
type SignUpPayload struct {
	Login    string `json:"login" validate:"required,min=6,max=30"`
	Name     string `json:"name" validate:"required,min=2,max=30"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

//easyjson:json
type AuthResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
}

//easyjson:skip
type UserClaims struct {
	ID   uuid.UUID `json:"-"`
	Role Role      `json:"-"`
	IP   string    `json:"-"`
}

//easyjson:json
type UpdatePasswordPayload struct {
	ID          uuid.UUID `json:"-"`
	OldPassword string    `json:"old_password"`
	NewPassword string    `json:"new_password"`
}

//easyjson:json
type UpdateRolePayload struct {
	ID   uuid.UUID `json:"-"`
	Role Role      `json:"role"`
}
