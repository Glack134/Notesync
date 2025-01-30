package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt        = "6765fgvbhhgf35vfu9jft5tg"
	signingKey  = "qjvkvnsjdnj2njn29njv**@9un19@!33"
	resetingKey = "fa#dh$bsia1*&2rffvsv2135v#eg*#"
	tokenTTL    = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user notesync.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) UpdatePasswordUser(username, password string) (string, error) {
	user, err := s.repo.UpdatePasswordUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) UpdatePasswordUserToken(token, newPassword string) error {
	isUsed, err := s.repo.IsTokenUsed(token)
	if err != nil {
		return err
	}
	if isUsed {
		return errors.New("token has already been used")
	}

	userID, err := s.repo.GetUserIDByToken(token)
	if err != nil {
		return err
	}

	newPasswordHash := s.generatePasswordHash(newPassword)

	err = s.repo.UpdatePasswordUserByID(userID, newPasswordHash)
	if err != nil {
		return err
	}

	err = s.repo.MarkTokenAsUsed(token)
	if err != nil {
		return err
	}

	return nil
}
