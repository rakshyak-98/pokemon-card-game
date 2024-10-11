package repository

import (
	"database/sql"
	"errors"

	"github.com/rakshyak-98/pokemonapi/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserById(id string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err := r.db.Exec(query, user.Username, user.Password)
	return err
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}

	query := "SELECT id, username, email FROM users WHERE username = ?"

	err := r.db.QueryRow(query, username).Scan(&user.Id, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserById(id string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, email FROM users WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&user.Id, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
	}
	return user, nil
}
