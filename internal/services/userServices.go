package services

import (
	"log"
	"net/http"

	// "strconv"

	"github.com/thepralad/socialnetwork/internal/models"
)

type UserService struct {
	userInteraction models.UserInteraction
}

func NewUserService(store models.UserInteraction) *UserService {
	return &UserService{userInteraction: store}
}

func (store *UserService) Poke(req *http.Request) string {
	token, _ := req.Cookie("session_token")
	if token.Value == "" {
		return "no token found"
	}

	sEmail, err := store.userInteraction.GetEmailFromToken(token.Value)
	if err != nil {
		log.Print(err)
	}

	rEmail := req.FormValue("receiver")

	store.userInteraction.Poke(sEmail, rEmail)

	return "Poked!"
}

func (store *UserService) AuthorizeUser(token string) (string, error) {
	email, err := store.userInteraction.GetEmailFromToken(token)
	if err != nil {
		log.Printf("Failed to get email from token: %v", err)
		return "", err
	}
	return email, nil
}

func (store *UserService) GetPokesService(email string) ([]models.Poke, error) {
	return store.userInteraction.GetPokes(email)
}
