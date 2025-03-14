package services

import (
	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/pkg/security"
)

type UserService struct{
	userStore models.UserStore
}

func NewUserService(store models.UserStore) *UserService{
	return &UserService{userStore: store}
}

func (s *UserService) RegisterUser(email, password string) error{
	hashedPassword, _ := security.HashPassword(password)
	s.userStore.CreateUser(email, hashedPassword)
	return nil
}
