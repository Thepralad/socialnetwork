package security

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), nil
}

func ComparePassword(password, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil	
}
