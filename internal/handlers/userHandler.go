package handlers

import(
	"net/http"
	"log"
	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/internal/services"
	"github.com/thepralad/socialnetwork/pkg/render"
)

type UserHandler struct{
	userService *services.UserService
}

func NewUserHandler(store models.UserStore) *UserHandler{
	userService := services.NewUserService(store)
	return &UserHandler{userService: userService}
}

func (h *UserHandler) HomeHandler(res http.ResponseWriter, req *http.Request){
	http.Redirect(res, req, "/register", http.StatusFound)
}

func (h *UserHandler) RegisterUser(res http.ResponseWriter, req *http.Request){
	//GET
	if req.Method == http.MethodGet{
		render.Template(res, "index", nil)
		return
	}

	//POST
	if req.Method == http.MethodPost{
		email := req.FormValue("email")
		password := req.FormValue("password")
		
		message, err := h.userService.RegisterUser(email, password)
		if err != nil{
			log.Print(err)
		}

		res.Write([]byte(message))
	}
}


