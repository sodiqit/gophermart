package dtos

import "time"

type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserUpdate struct {
	Login        *string
	PasswordHash *string
}
