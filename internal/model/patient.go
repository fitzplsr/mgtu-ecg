package model

import (
	"time"
)

const BirthdayFormat = "1990-05-15"

//easyjson:json
type Patient struct {
	ID        int       `json:"int"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//easyjson:json
type Patients struct {
	Patients []*Patient `json:"patients"`
}

//easyjson:json
type CreatePatient struct {
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Birthday time.Time `json:"birthday"`
}

//easyjson:json
type GetPatient struct {
	ID int `json:"id"`
}
