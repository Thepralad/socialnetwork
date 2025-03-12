package models

import (
	"database/sql"
	"time"
)

type MySQLUserStore struct{
	DB *sql.DB
}

func (m *MySQLUserStore) CreateUser(email string, password string) (int64, error){
	query := "INSERT INTO users(email, password, created_at) VALUES(?,?,?)"
	result, err := m.DB.Exec(query, email, password, time.Now())
	if err != nil{
		return 0, err
	}
	return result.LastInsertId()
}

func (m *MySQLUserStore) GetUserByEmail(email string) (*User, error){
	var user User
	query := "SELECT id, email, password, created_at FROM users where email = ?"
	err := m.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt) 
	if err != nil{
		return nil, err
	}
	return &user, nil
}
