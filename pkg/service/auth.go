package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/repository"
)

const (
	salt       = "6765fgvbhhgf35vfu9jft5tg"
	signingKey = "qjvkvnsjdnj2njn29njv**@9un19@!33"
	tokenTTL   = 12 * time.Hour
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

func (s *AuthService) SendPasswordResetEmail(email string) error {
	// Проверка на пустой email
	if email == "" {
		return errors.New("email cannot be empty")
	}
	userID := 21
	// Проверка существования пользователя
	user, err := s.repo.FindUserByEmail(userID, email) // Предположим, у вас есть функция для поиска пользователя по email
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == (notesync.User{}) {
		return errors.New("user not found")
	}

	// Настройки SMTP (пример для Gmail)
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"
	sender := "rbhb05@mail.ru"         // Замените на ваш адрес электронной почты
	password := "Lvpbhiyaw6i9KGaWgePj" // Замените на ваш пароль или используйте App Password

	// Создаем сообщение
	subject := "Subject: Password Reset Request\n"
	body := "Please click the link to reset your password: http://localhost:8080/auth/reset-password\n"
	message := []byte(subject + "\n" + body)

	// Отправляем email
	auth := smtp.PlainAuth("", sender, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{email}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
