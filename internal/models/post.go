package models

type Post struct {
	ID        int64
	Email     string
	Content   string
	CreatedAt string
	Metric    int64
}

type PostStore interface {
	GetEmailFromToken(token string) (string, error)
}

