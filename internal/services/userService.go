package services

import (
	"fmt"
)

func RegisterUser(email string, password string) string{
	return fmt.Sprintf("User registered as: %s | %s", email, password)
}
