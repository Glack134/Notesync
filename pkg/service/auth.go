package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt        = "6765fgvbhhgf35vfu9jft5tg"
	signingKey  = "qjvkvnsjdnj2njn29njv**@9un19@!33"
	resetingKey = "fa#dh$bsia1*&2rffvsv2135v#eg*#"
	tokenTTL    = 12 * time.Hour
	tokenEmail  = 10 * time.Minute
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

func GenerateTokenEmail() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *AuthService) CreateResetToken(email string) (string, error) {
	userID, err := s.repo.GetTokenResetPassword(email)
	if err != nil {
		return "", err
	}

	token, err := GenerateTokenEmail()
	if err != nil {
		return "", err
	}

	expiry := time.Now().Add(tokenEmail)
	err = s.repo.SaveResetToken(userID, token, expiry)
	if err != nil {
		return "", err
	}

	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)

	// Теперь отправляем email с ссылкой для сброса пароля
	subject := "Сброс пароля"
	body := fmt.Sprintf("Вы запросили сброс пароля. Пожалуйста, перейдите по следующей ссылке, чтобы сбросить ваш пароль: %s", resetLink)

	if err := s.sendEmail(email, subject, body); err != nil {
		return "", err
	}

	return resetLink, nil
}

func (s *AuthService) sendEmail(to string, subject string, body string) error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Ошибка загрузки данных фаил: %s", err)
		return err
	}

	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	// Настройка SMTP-сервера
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

	// Подготовка сообщения
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// Аутентификация
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправка письма
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Printf("Ошибка при отправке email: %s", err)
		return err
	}

	log.Printf("Инструкция по восстановлению отправлена на почту %s", to)
	return nil
}
