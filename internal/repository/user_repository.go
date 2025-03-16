package repository

import (
	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user domain.User) error {
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, user.Username, user.Email, user.Password)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email, password FROM users WHERE username = ?`
	row := r.DB.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}