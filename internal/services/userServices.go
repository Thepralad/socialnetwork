package services

import (
	// "log"
	// "net/http"
	// "strconv"

	"github.com/thepralad/socialnetwork/internal/models"
)

type UserService struct {
	userStore models.UserInteraction
}

func NewUserService(store models.UserInteraction) *UserService {
	return &UserService{userStore: store}
}

func (store *UserService) Poke(to, from string){
	
}