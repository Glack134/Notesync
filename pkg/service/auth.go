package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/repository"
)

const salt = "6765fgvbhhgf35vfu9jft5tg"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user notesync.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
