package models

import (
	"database/sql"
	"time"
)

type MySQLUserStore struct{
	DB *sql.DB
}

//-----USER-----
func (m *MySQLUserStore) CreateUser(email string, password string, token string) (int64, error){
	query := "INSERT INTO users(email, password, created_at, token) VALUES(?,?,?,?)"
	result, err := m.DB.Exec(query, email, password, time.Now(), token)
	if err != nil{
		return 0, err
	}
	return result.LastInsertId()
}

func (m *MySQLUserStore) GetUserByEmail(email string) (*User, error){
	var user User
	query := "SELECT id, email, password, created_at, verified FROM users where email = ?"
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.Verified) 
	if err != nil{
		return nil, err
	}
	return &user, nil
}
func (m *MySQLUserStore) SetSessionToken(email, token string) error{
	query := "INSERT INTO sessions (email, token) VALUES (?, ?)"
	_, err := m.DB.Exec(query, email, token)
	if err != nil{
		return err
	}
	return nil
}

func (m *MySQLUserStore) VerifyUserByToken(token string) error{
	query := "UPDATE users SET verified = 1 WHERE token = ?"
	_, err := m.DB.Exec(query, token)
	if err != nil{
		return err
	}
	return nil
}

//-----POSTS-----
func (m *MySQLUserStore) GetEmailFromToken(token string) (string, error){
	var email string
	query := "SELECT email FROM sessions WHERE token = ?"
	err := m.DB.QueryRow(query, token).Scan(&email)
	if err != nil{
		return "", err
	}
	return email, err
}
