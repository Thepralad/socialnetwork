package services

import (
	"log"
	"net/http"

	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/pkg/mail"
	"github.com/thepralad/socialnetwork/pkg/security"
	"github.com/thepralad/socialnetwork/pkg/token"
)

type AuthService struct {
	userStore models.UserStore
}

func NewAuthService(store models.UserStore) *AuthService {
	return &AuthService{userStore: store}
}

// Business logic for Registering user - Returns the status message
func (s *AuthService) RegisterUser(email, password string) (string, error) {

	isValidEmail := security.IsSalesianEmail(email)
	if !isValidEmail {
		return "Only @salesiancollege.net emails hehe", nil
	}

	//Checks if user already exists
	userAlreadyExist, err := s.userStore.GetUserByEmail(email)
	if userAlreadyExist != nil {
		return "User already exists", err
	}

	//Ensuring the password is or more than 8 character
	if len(password) < 8 {
		return "I guess the password is not so strong...", nil
	}

	//Encrypting password using bcrypt
	hashedPassword, _ := security.HashPassword(password)
	verificationToken := token.GenerateToken()
	//If no above condition is met, the email and password in inserted to the DB
	s.userStore.CreateUser(email, hashedPassword, verificationToken)
	message := "http://snet.club/verify?token=" + verificationToken

	err = mail.SendEmail([]string{email}, "Welcome to SNET!", message)
	if err != nil {
		log.Fatal(err)
	}
	err = s.userStore.CreateUserProfile(email)
	if err != nil {
		log.Print("failed creating profile")
	}
	return "Check your email", nil
}

// Business logic for Logging in user - Returns the status message
func (s *AuthService) LoginUser(email, password string) (string, error) {
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

func (s *AuthService) AddSessionToken(email string, res http.ResponseWriter) {
	sessionToken := token.GenerateToken()

	// Set Token in the DB with the email address
	err := s.userStore.SetSessionToken(email, sessionToken)
	if err != nil {
		log.Printf("Failed to set session token: %v", err)
		return
	}

	//Set Cookie with proper attributes
	http.SetCookie(res, &http.Cookie{
		Name:   "session_token",
		Value:  sessionToken,
		MaxAge: 2592000, // 30 days
	})
}

func (s *AuthService) RemoveSessionToken(res http.ResponseWriter) {
	http.SetCookie(res, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: 3600,
	})
}

// Business logic for Verifying user - Returns the status message
func (s *AuthService) VerifyUser(token string) error {
	err := s.userStore.VerifyUserByToken(token)
	if err != nil {
		return err
	}
	return nil
}
