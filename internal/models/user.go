package models

type User struct {
	ID        int64
	Email     string
	Password  string
	CreatedAt string
	Verified  bool
	Token     string
}
type UserProfile struct {
	Email              string
	Username           string
	InstagramLink      string
	Gender             string
	TopArtist          string
	RelationshipStatus string
	LookingFor         string
	FactAboutMe        string
	Dept               string
	Year               string
	ImgURL             string
}

type Poke struct {
	Username string
	Email    string
	ImgURL   string
}

type UserStore interface {
	CreateUser(email string, password string, token string) (int64, error)
	CreateUserProfile(email string) error
	GetUserByEmail(email string) (*User, error)
	VerifyUserByToken(token string) error
	SetSessionToken(email, token string) error
}

type UserInteraction interface {
	Poke(to, from string)
	GetEmailFromToken(token string) (string, error)
	GetPokes(email string) ([]Poke, error)
}
