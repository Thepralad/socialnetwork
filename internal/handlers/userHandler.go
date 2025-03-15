package handlers

import(
	"net/http"
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
	if req.Method == http.MethodGet{
		render.Template(res, "index", nil)
		return
	}

	if req.Method == http.MethodPost{

		email := req.FormValue("email")
		password := req.FormValue("password")
		
		h.userService.RegisterUser(email, password)

		res.Write([]byte("check your email"))
	}
}


