package model

import (
	"github.com/google/uuid"
	"time"
)

//easyjson:json
type Session struct {
	UserId       uuid.UUID     `json:"user_id"`
	UserRole     Role          `json:"user_role"`
	RefreshToken string        `json:"refresh_token"`
	UserAgent    string        `json:"user_agent"`
	Fingerprint  string        `json:"fingerprint"`
	Ip           string        `json:"ip"`
	ExpiresIn    time.Duration `json:"expires_in"`
	CreatedAt    time.Time     `json:"created_at"`
}
