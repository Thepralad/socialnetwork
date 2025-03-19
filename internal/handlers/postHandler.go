package handlers

import (
	"log"
	"net/http"

	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/internal/services"
	"github.com/thepralad/socialnetwork/pkg/render"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(store models.PostStore) *PostHandler {
	postService := services.NewPostService(store)
	return &PostHandler{postService: postService}
}

func (h *PostHandler) HomePostHandler(res http.ResponseWriter, req *http.Request) {
	// Extract session_token from request
	if req.Method == http.MethodGet {
		token, _ := req.Cookie("session_token")
		if token.Value == "" {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		email, _ := h.postService.AuthorizeUser(token.Value)
		data := map[string]interface{}{
			"Email": email,
		}
		if err := render.Template(res, "feeds", data); err != nil {
			http.Error(res, "Error rendering template", http.StatusInternalServerError)
			log.Printf("Template error: %v", err)
			return
		}

	}
}
