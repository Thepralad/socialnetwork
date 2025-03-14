package handlers

import(
	"net/http"
	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/internal/services"
)

type UserHandler struct{
	userService *services.UserService
}

func NewUserHandler(store models.UserStore) *UserHandler{
	userService := services.NewUserService(store)
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser(res http.ResponseWriter, req *http.Request){
	//Getting user input from the form
	email := req.FormValue("email")
	password := req.FormValue("password")
	
	h.userService.RegisterUser(email, password)

	res.Write([]byte("user created"))
}


