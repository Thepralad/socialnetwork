package services

import (
	"log"
	"net/http"
	"strconv"

	"github.com/thepralad/socialnetwork/internal/models"
)

type PostService struct {
	postStore models.PostStore
}

func NewPostService(store models.PostStore) *PostService {
	return &PostService{postStore: store}
}

func (store *PostService) CreatePost(post *models.Post) error {
	log.Println("Just creates a simple post")
	return nil
}

func (store *PostService) AuthorizeUser(token string) (string, error) {
	email, err := store.postStore.GetEmailFromToken(token)
	if err != nil {
		log.Printf("Failed to get email from token: %v", err)
		return "", err
	}
	return email, nil
}

func (store *PostService) Post(email, content string) error {
	err := store.postStore.Post(email, content)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil

}

func (store *PostService) UpdateMetric(req *http.Request) int {
	postIdStr := req.URL.Query().Get("postid")
	action := req.URL.Query().Get("action")

	postIdInt, _ := strconv.Atoi(postIdStr)
	err := store.postStore.MetricUpdate(action, postIdInt)
	if err != nil {
		log.Print(err)
	}

	count, err := store.postStore.GetMetric(postIdInt)
	return count
}

func (store *PostService) GetPostsService(offset int) ([]models.Post, error) {
	return store.postStore.GetPosts(offset)
}

func (store *PostService) GetPostsByUserService(email string) ([]models.Post, error) {
	return store.postStore.GetPostsByUser(email)
}

func (store *PostService) GetPokesService(email string) ([]models.Poke, error) {
	return store.postStore.GetPokes(email)
}

func (store *PostService) GetProfile(email string) (*models.UserProfile, error) {
	profile, err := store.postStore.GetProfileFromEmail(email)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return profile, nil

}

func (store *PostService) UpdateProfile(email string, profile *models.UserProfile) error {
	err := store.postStore.UpdateProfileFromEmail(email, profile)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
