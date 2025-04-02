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
	token, _ := req.Cookie("session_token")
	if token.Value == "" {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	email, _ := h.userService.AuthorizeUser(token.Value)
	t, _ := template.New("pokes").Parse(`{{range .Pokes}}
							<div class="poke-item">
								<img src="{{.ImgURL}}" alt="">
								<div class="poke-info">
									<p>@{{.Username}}</p>
									<p><a href="#">{{.Email}}</a></p>
								</div>
								<form hx-post="/pokeback?email={{.Email}}" hx-target="this" hx-swap="outerHTML">
									<button type="submit" class="poke-back-btn">Poke Back</button>
								</form>
							</div>
						{{end}}`)
	pokes, _ := h.userService.GetPokesService(email)
	data := struct {
		Pokes []models.Poke
	}{
		Pokes: pokes,
	}
	err := t.Execute(res, data)
	if err != nil {
		log.Printf("Template error: %v", err)
		return
	}
}
