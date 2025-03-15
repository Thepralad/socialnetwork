package services

import (
	"log"

	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/pkg/mail"
	"github.com/thepralad/socialnetwork/pkg/security"
	"github.com/thepralad/socialnetwork/pkg/token"
)

type UserService struct {
	userStore models.UserStore
}

func NewUserService(store models.UserStore) *UserService {
	return &UserService{userStore: store}
}

// Business logic for Registering user - Returns the status message
func (s *UserService) RegisterUser(email, password string) (string, error) {
	//Checks if user already exists
	userAlreadyExist, err := s.userStore.GetUserByEmail(email)
	if userAlreadyExist != nil {
		return "User already exists", err
	}

	//Ensuring the password is or more than 8 character
	if len(password) < 8 {
		return "Password atleast 8 characters", nil
	}

	//Encrypting password using bcrypt
	hashedPassword, _ := security.HashPassword(password)
	verificationToken := token.GenerateVerificationToken()
	//If no above condition is met, the email and password in inserted to the DB
	s.userStore.CreateUser(email, hashedPassword, verificationToken)
	message := "http://localhost:8080/verify?token=" + verificationToken

	err = mail.SendEmail([]string{email}, "Welcome to SNET!", message)
	if err != nil {
		log.Fatal(err)
	}
	return "Check your email", nil
}

// Business logic for Logging in user - Returns the status message
func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.userStore.GetUserByEmail(email)
	if err != nil {
		return "User does not exists", nil
	}

	if !user.Verified {
		return "User not verified", nil
	}

	if !security.ComparePassword(password, user.Password) {
		return "Password does not match", nil
	}
	return "user logged in successfully", nil
}

// Business logic for Verifying user - Returns the status message
func (s *UserService) VerifyUser(token string) error {
	err := s.userStore.VerifyUserByToken(token)
	if err != nil {
		return err
	}
	return nil
}
