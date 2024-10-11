package services

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rakshyak-98/pokemonapi/internal/models"
	"github.com/rakshyak-98/pokemonapi/internal/repository"
)

type UserService interface {
	RegisterUser(user *models.User) error
	AuthenticateUser(username, password string) (string, error)
	GetUser(id string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) RegisterUser(user *models.User) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

func (s *userService) AuthenticateUser(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if !comparePasswords(user.Password, password) {
		return "", errors.New("invalid password")
	}

	token, err := generateJWTToken(user.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUser(id string) (*models.User, error) {
	return s.repo.GetUserById(id)
}

func hashPassword(password string) (string, error) {
	hashBuffer := sha256.New()
	hashBuffer.Write([]byte(password))

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hashBuffer.Write(salt)
	hashed := hashBuffer.Sum(nil)
	return base64.StdEncoding.EncodeToString(salt) + ":" + base64.StdEncoding.EncodeToString(hashed), nil
}

func comparePasswords(hashedPassword, password string) bool {
	parts := strings.Split(hashedPassword, ":")
	if len(parts) != 2 {
		return false // invalid has format
	}

	salt, _ := base64.StdEncoding.DecodeString(parts[0])
	pass, _ := base64.StdEncoding.DecodeString(parts[1])

	hashBuffer := sha256.New()
	hashBuffer.Write(salt)
	hashBuffer.Write([]byte(password))

	return bytes.Equal(hashBuffer.Sum(nil), pass)
}

func generateJWTToken(userID int, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(secret)
}
