package models


type User struct {
	ID int64
	Email string
	Password string
	CreatedAt string
	Verified bool
	Token string
}

type UserStore interface {
	CreateUser(email string, password string, token string) (int64, error)
	GetUserByEmail(email string) (*User, error)
	VerifyUserByToken(token string) error
}
