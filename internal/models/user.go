package models

import (
	"time"
)

type User struct {
	ID int64
	Email string
	Password string
	CreatedAt time.Time
}

type UserStore interface {
	CreateUser(email string, password string) (int64, error)
	GetUserByEmail(email string) (*User, error)
}
