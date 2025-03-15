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

//Business logic for Registering user - Returns the status message
func (s *UserService) RegisterUser(email, password string) (string, error){
	//Checks if user already exists
	userAlreadyExist, err := s.userStore.GetUserByEmail(email)
	if userAlreadyExist != nil {
		return "User already exists", err 
	}

	//Ensuring the password is or more than 8 character
	if len(password) < 8{
		return "Password atleast 8 characters", nil
	}

	//Encrypting password using bcrypt
	hashedPassword, _ := security.HashPassword(password)

	//If no above condition is met, the email and password in inserted to the DB
	s.userStore.CreateUser(email, hashedPassword)
	
	return "Check your email", nil
}

func (s *UserService) GetUser(email string) (*models.User, error){
	user, err := s.userStore.GetUserByEmail(email)
	if err != nil{
		return nil, err
	}
	return user, nil
}
