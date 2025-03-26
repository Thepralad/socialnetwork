package handlers

import (
	"fmt"
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
		profile, err := h.postService.GetProfile(email)
		if err != nil {
			log.Print(err)
		}
		posts, _ := h.postService.GetPostsService(0) // Get initial posts
		data := map[string]interface{}{
			"Profile": profile,
			"Posts":   posts,
		}

		if err := render.Template(res, "feeds", data); err != nil {
			http.Error(res, "Error rendering template", http.StatusInternalServerError)
			log.Printf("Template error: %v", err)
			return
		}

	}
}


func (h *PostHandler) EditProfileHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		token, _ := req.Cookie("session_token")
		if token.Value == "" {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		email, _ := h.postService.AuthorizeUser(token.Value)
		profile, _ := h.postService.GetProfile(email)
		data := map[string]interface{}{
			"Profile": profile,
		}

		if err := render.Template(res, "editprofile", data); err != nil {
			http.Error(res, "Error rendering template", http.StatusInternalServerError)
			log.Printf("Template error: %v", err)
			return
		}
	}

	if req.Method == http.MethodPost {
		token, _ := req.Cookie("session_token")
		if token.Value == "" {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

		email, _ := h.postService.AuthorizeUser(token.Value)
		username := req.FormValue("username")
		gender := req.FormValue("gender")
		relationship_status := req.FormValue("relationship_status")
		top_artist := req.FormValue("top_artist")
		looking_for := req.FormValue("looking_for")
		fact_about_me := req.FormValue("fact_about_me")
		dept := req.FormValue("dept")
		year := req.FormValue("year")

		profile := models.UserProfile{
			Email:              email,
			Username:           username,
			Gender:             gender,
			RelationshipStatus: relationship_status,
			TopArtist:          top_artist,
			LookingFor:         looking_for,
			FactAboutMe:        fact_about_me,
			Dept:				dept,
			Year: 				year,
		}

		err := h.postService.UpdateProfile(email, &profile)
		if err != nil {
			log.Printf("Error updating profile: %v", err)
			// Optionally return an error message to frontend
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(res, req, "/home", http.StatusSeeOther)
	}
}

func (h *PostHandler) PostHandler(res http.ResponseWriter, req *http.Request){
	token, _ := req.Cookie("session_token")
		if token.Value == "" {
			http.Redirect(res, req, "/login", http.StatusFound)
			return
		}

	email, _ := h.postService.AuthorizeUser(token.Value)
	content := req.FormValue("content")
	h.postService.Post(email, content)
	
		http.Redirect(res, req, "/home", http.StatusSeeOther)
}

func (h *PostHandler) GetPostHandler(res http.ResponseWriter, req *http.Request) {
	offsetStr := req.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		fmt.Sscanf(offsetStr, "%d", &offset)
	}

	posts, err := h.postService.GetPostsService(offset)
	if err != nil {
		http.Error(res, "Error fetching posts", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	data := map[string]interface{}{
		"Posts": posts,
	}

	if err := render.Template(res, "post_item", data); err != nil {
		http.Error(res, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}
}
