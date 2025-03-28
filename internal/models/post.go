package models


type Post struct {
    ID       int
    Username string
    ImgURL string
    Email    string
    Content  string
    VoteCount int
}
type PostStore interface {
	GetEmailFromToken(token string) (string, error)
	GetProfileFromEmail(email string) (*UserProfile, error)
	UpdateProfileFromEmail(email string, profile *UserProfile) error
    GetPosts(offset int) ([]Post, error)
    Post(email,content string) error
    MetricUpdate(action, post_id int) error
}


