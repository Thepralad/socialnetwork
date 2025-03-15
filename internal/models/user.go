package models


type User struct {
	ID int64
	Email string
	Password string
	CreatedAt string
}

type UserStore interface {
	CreateUser(email string, password string) (int64, error)
	GetUserByEmail(email string) (*User, error)
}
