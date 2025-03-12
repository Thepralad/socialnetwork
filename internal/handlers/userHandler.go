package handlers

import(
	"net/http"
	"github.com/thepralad/socialnetwork/internal/models"
)

type UserHandler struct{
	userStorage models.UserStore
}

func NewUserHandler(store models.UserStore) *UserHandler{
	return &UserHandler{userStorage: store}
}

func (h *UserHandler) RegisterUser(res http.ResponseWriter, req *http.Request){
	//Getting user input from the form
	email := req.FormValue("email")
	password := req.FormValue("password")
	
	h.userStorage.CreateUser(email, password)

	res.Write([]byte("user created"))
}


