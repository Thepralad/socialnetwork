package handlers

import (
	"log"
	"net/http"

	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/internal/services"
	"github.com/thepralad/socialnetwork/pkg/render"
)

type AuthHandler struct {
	userService *services.AuthService
}

func NewAuthHandler(store models.UserStore) *AuthHandler {
	userService := services.NewAuthService(store)
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) HomeHandler(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "/register", http.StatusFound)
}

func (h *AuthHandler) RegisterHandler(res http.ResponseWriter, req *http.Request) {
	//GET
	if req.Method == http.MethodGet {
		token, _ := req.Cookie("session_token")
		if token.Value == "" {
			if err := render.Template(res, "index", nil); err != nil {
				log.Printf("Template error: %v", err)
				return
			}
			return
		}
		http.Redirect(res, req, "/home", http.StatusFound)
		return
	}

	//POST
	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")

		message, err := h.userService.RegisterUser(email, password)
		if err != nil {
			log.Printf("Registration error: %v", err)
			return
		}

		res.Write([]byte(message))
	}
}

func (h *AuthHandler) LoginHandler(res http.ResponseWriter, req *http.Request) {
	//GET
	if req.Method == http.MethodGet {
		token, _ := req.Cookie("session_token")
		if token.Value == "" {
			if err := render.Template(res, "login", nil); err != nil {
				log.Printf("Template error: %v", err)
				return
			}
			return
		}

		http.Redirect(res, req, "/home", http.StatusFound)
		return
	}

	//POST
	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")

		message, err := h.userService.LoginUser(email, password)
		if err != nil {
			log.Printf("Login error: %v", err)
			return
		}

		// Check login message for success
		if message != "user logged in successfully" {
			// Return to login page with error
			data := map[string]interface{}{
				"Error": message,
			}
			if err := render.Template(res, "login", data); err != nil {
				log.Printf("Template error: %v", err)
			}
			return
		}

		h.userService.AddSessionToken(email, res)
		// Successful login - redirect to home
		http.Redirect(res, req, "/home", http.StatusSeeOther)
		return
	}
}

func (h *AuthHandler) LogoutHandler(res http.ResponseWriter, req *http.Request) {
	h.userService.RemoveSessionToken(res)
	http.Redirect(res, req, "/login", http.StatusFound)
}

func (h *AuthHandler) VerifyHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	token := req.URL.Query().Get("token")

	if err := h.userService.VerifyUser(token); err != nil {
		log.Printf("Verification error: %v", err)
		return
	}

	if err := render.Template(res, "verified", nil); err != nil {
		log.Printf("Template error: %v", err)
		return
	}
}
