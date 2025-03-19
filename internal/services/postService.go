package services

import (
	"log"

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
