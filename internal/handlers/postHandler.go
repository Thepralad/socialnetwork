package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/thepralad/socialnetwork/internal/models"
	"github.com/thepralad/socialnetwork/internal/services"
	"github.com/thepralad/socialnetwork/pkg/cloudstorage"
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

		imgFile, _, _ := req.FormFile("profile_img")

		img_url, _ := cloudstorage.UploadImg(imgFile)

		profile := models.UserProfile{
			Email:              email,
			Username:           username,
			Gender:             gender,
			RelationshipStatus: relationship_status,
			TopArtist:          top_artist,
			LookingFor:         looking_for,
			FactAboutMe:        fact_about_me,
			Dept:               dept,
			Year:               year,
			ImgURL:             img_url,
		}

		if img_url == "" {
			existingUserInfo, _ := h.postService.GetProfile(email)
			profile.ImgURL = existingUserInfo.ImgURL
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

func (h *PostHandler) PostHandler(res http.ResponseWriter, req *http.Request) {
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

	res.Header().Set("Content-Type", "text/html")

	for _, post := range posts {
		fmt.Fprintf(res, `
			<div class="post">
				<div class="upper-side">
					<img src="%s" class="post-profile-pic">
					<div>
						<p class="post-username">@%s</p>
						<p class="post-email"><a href="http://localhost:8080/home/%s">%s</a></p>
					</div>
				</div>
				<p class="post-text">%s</p>
				<div class="actionBtns">
					<button class="upvote" hx-get="/updatemetric?postid=%d&action=up" hx-target="#vote-count-%d" hx-swap="innerHTML">⇧</button>
					<p id="vote-count-%d">%d</p>
					<button class="downvote" hx-get="/updatemetric?postid=%d&action=down" hx-target="#vote-count-%d" hx-swap="innerHTML">⇩</button>
				</div>
				<p class="post-time">Just now</p>
			</div>
			
		`, post.ImgURL, post.Username, post.Email, post.Email, post.Content, post.ID, post.ID, post.ID, post.VoteCount, post.ID, post.ID)
	}

}



func (h *PostHandler) UserProfileHandler(res http.ResponseWriter, req *http.Request) {

	email := strings.TrimPrefix(req.URL.Path, "/home/")

	profile, _ := h.postService.GetProfile(email)

	data := map[string]interface{}{
		"Profile": profile,
	}

	if err := render.Template(res, "userprofile", data); err != nil {
		http.Error(res, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}
}

func (h *PostHandler) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	count := h.postService.UpdateMetric(req)
	fmt.Fprint(res, count)
}
