package handlers

import (
	"net/http"
	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/internal/services"
)

type UserHandler struct{
	userService *services.UserService
}

func NewUserHandler(store models.UserInteraction) *UserHandler {
	userService := services.NewUserService(store)
	return &UserHandler{userService: userService}
}

func (h *UserHandler)PokeHandler(res http.ResponseWriter, req *http.Request){
	res.Write([]byte("this thing sucks"))
}