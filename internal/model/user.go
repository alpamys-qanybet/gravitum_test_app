package model

import "time"

type CreateUserRequest struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
}

type UpdateUserRequest struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
}

type User struct {
	Id         uint       `json:"id"`
	Name       string     `json:"name"`
	Surname    *string    `json:"surname,omitempty"`
	InsertedAt time.Time  `json:"inserted_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}
