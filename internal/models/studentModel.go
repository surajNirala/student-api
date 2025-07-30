package models

import "time"

type Student struct {
	Id        uint64    `json:"id"`
	Name      string    `validate:"required" json:"name"`
	Email     string    `validate:"required" json:"email"`
	Age       int       `validate:"required" json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
