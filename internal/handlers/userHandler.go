package handlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(store models.UserInteraction) *UserHandler {
	userService := services.NewUserService(store)
	return &UserHandler{userService: userService}
}

func (h *UserHandler) PokeHandler(res http.ResponseWriter, req *http.Request) {
	message := h.userService.Poke(req)
	if message == "no token found" {
		http.Redirect(res, req, "/home", http.StatusFound)
	}

	res.Write([]byte(message))
}

func (h *UserHandler) PokesHandler(res http.ResponseWriter, req *http.Request) {
	token, err := req.Cookie("session_token")
	if err != nil || token.Value == "" {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	email, err := h.userService.AuthorizeUser(token.Value)
	if err != nil {
		log.Printf("Error authorizing user: %v", err)
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	t, err := template.New("pokes").Parse(`{{range .Pokes}}
							<div class="poke-item">
								<img src="{{.ImgURL}}" alt="">
								<div class="poke-info">
									<p><a href="http://snet.club/home/{{.Email}}">@{{.Username}}</a></p>
								</div>
								<form hx-post="/pokeback?email={{.Email}}" hx-target="this" hx-swap="outerHTML">
									<button type="submit" class="poke-back-btn">Poke</button>
								</form>
							</div>
						{{end}}`)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	pokes, err := h.userService.GetPokesService(email)
	if err != nil {
		log.Printf("Error getting pokes: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Pokes []models.Poke
	}{
		Pokes: pokes,
	}

	if err := t.Execute(res, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
}
