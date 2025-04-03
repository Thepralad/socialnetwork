package models

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

type MySQLUserStore struct {
	DB *sql.DB
}

// -----USER-----

// Add a user (email, password, token) to database. [snet.users]
func (m *MySQLUserStore) CreateUser(email string, password string, token string) (int64, error) {
	query := "INSERT INTO users(email, password, created_at, token) VALUES(?,?,?,?)"
	result, err := m.DB.Exec(query, email, password, time.Now(), token)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Get *User by using email [snet.users]
func (m *MySQLUserStore) GetUserByEmail(email string) (*User, error) {
	var user User
	query := "SELECT id, email, password, created_at, verified FROM users where email = ?"
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.Verified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Associates a session token with email [snet.sessions]
func (m *MySQLUserStore) SetSessionToken(email, token string) error {
	query := "INSERT INTO sessions (email, token) VALUES (?, ?)"
	_, err := m.DB.Exec(query, email, token)
	if err != nil {
		return err
	}
	return nil
}

// Updates the `snet.users.verified` to 1
func (m *MySQLUserStore) VerifyUserByToken(token string) error {
	query := "UPDATE users SET verified = 1 WHERE token = ?"
	_, err := m.DB.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}

// Creates a new row in the [snet.user_profiles] with inital values (username, email)
func (m *MySQLUserStore) CreateUserProfile(email string) error {
	query := "INSERT INTO user_profiles (username, email) VALUES(?, ?)"
	_, err := m.DB.Exec(query, extractUsername(email), email)
	if err != nil {
		return err
	}
	log.Printf("User profle created %s", extractUsername(email))
	return nil
}

// Extracts the username from the email
func extractUsername(email string) string {
	atIndex := strings.Index(email, "@")
	return email[:atIndex]
}

// -----POSTS-----

// Inserts a new row in [snet.posts] with (email, content)
func (m *MySQLUserStore) Post(email, content string) error {
	query := "INSERT INTO posts(email, content) VALUES(?, ?)"
	_, err := m.DB.Exec(query, email, content)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves all the data from [snet.posts] & [snet.user_profiles]. (id, email, username, img_url, content, vote_count)
func (m *MySQLUserStore) GetPosts(offset int) ([]Post, error) {
	query := `
		SELECT 
	    p.id, 
	    p.email, 
	    COALESCE(u.username, 'Unknown') AS username, 
	    COALESCE(u.img_url, 'default.jpg') AS img_url, 
	    p.content, 
	    p.vote_count 
		FROM posts p
		LEFT JOIN user_profiles u ON p.email = u.email
		ORDER BY p.created_at DESC
		LIMIT 20 OFFSET ?;
		`
	rows, err := m.DB.Query(query, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Email, &p.Username, &p.ImgURL, &p.Content, &p.VoteCount)
		posts = append(posts, p)
	}

	return posts, nil
}

func (m *MySQLUserStore) GetPostsByUser(email string) ([]Post, error) {
	query := `
		SELECT 
		p.id, 
		p.email, 
		COALESCE(u.username, 'Unknown') AS username, 
		COALESCE(u.img_url, 'default.jpg') AS img_url, 
		p.content, 
		p.vote_count 
		FROM posts p
		LEFT JOIN user_profiles u ON p.email = u.email
		WHERE p.email = ?
		ORDER BY p.created_at DESC;
		`
	rows, err := m.DB.Query(query, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Email, &p.Username, &p.ImgURL, &p.Content, &p.VoteCount)
		posts = append(posts, p)
	}

	return posts, nil
}

// Retrieves Email using Token [snet.sessions]
func (m *MySQLUserStore) GetEmailFromToken(token string) (string, error) {
	var email string
	query := "SELECT email FROM sessions WHERE token = ?"
	err := m.DB.QueryRow(query, token).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, err
}

// Retrieves *UserProfile using email [snet.user_profiles]
func (m *MySQLUserStore) GetProfileFromEmail(email string) (*UserProfile, error) {
	var profile UserProfile
	query := `SELECT email, username, instagram_link, gender, top_artist, relationship_status, looking_for, fact_about_me, dept, year, img_url
        FROM user_profiles
        WHERE email = ?;`
	err := m.DB.QueryRow(query, email).Scan(&profile.Email,
		&profile.Username,
		&profile.InstagramLink,
		&profile.Gender,
		&profile.TopArtist,
		&profile.RelationshipStatus,
		&profile.LookingFor,
		&profile.FactAboutMe,
		&profile.Dept,
		&profile.Year,
		&profile.ImgURL,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Takes in action(up/down) & postid to update the metric
func (m *MySQLUserStore) MetricUpdate(action string, post_id int) error {
	query := ""

	if action == "up" {
		query = "UPDATE posts SET vote_count = vote_count + 1  WHERE id = ?;"
	} else if action == "down" {
		query = "UPDATE posts SET vote_count = vote_count - 1  WHERE id = ?;"
	}

	_, err := m.DB.Exec(query, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQLUserStore) GetMetric(post_id int) (int, error) {
	var count int

	query := "SELECT vote_count FROM posts WHERE id = ?"

	err := m.DB.QueryRow(query, post_id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *MySQLUserStore) UpdateProfileFromEmail(email string, profile *UserProfile) error {
	query := `
	    UPDATE user_profiles
	    SET 
	        username = ?, 
	        instagram_link = ?, 
	        gender = ?, 
	        top_artist = ?, 
	        relationship_status = ?, 
	        looking_for = ?, 
	        fact_about_me = ?,
	        dept = ?,
	        year = ?,
	        img_url = ?
	    WHERE email = ?;
	`

	_, err := m.DB.Exec(query, profile.Username, profile.InstagramLink, profile.Gender, profile.TopArtist, profile.RelationshipStatus, profile.LookingFor, profile.FactAboutMe, profile.Dept, profile.Year, profile.ImgURL, profile.Email)
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQLUserStore) GetPokes(email string) ([]Poke, error) {
	query := "SELECT u.username, u.img_url, u.email FROM pokes p JOIN user_profiles u ON p.sender = u.email WHERE p.receiver = ?;"

	rows, err := m.DB.Query(query, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pokes []Poke
	for rows.Next() {
		var p Poke
		err := rows.Scan(&p.Username, &p.ImgURL, &p.Email)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		pokes = append(pokes, p)
	}

	return pokes, nil
}

// USER

func (m *MySQLUserStore) Poke(to, from string) {
	_, err := m.DB.Exec("INSERT INTO pokes (sender, receiver) VALUES (?, ?)", to, from)
	if err != nil {
		log.Print(err)
	}
}
