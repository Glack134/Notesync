package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/polyk005/notesync/pkg/repository"
)

const (
	tokenEmail = 10 * time.Minute
)

type AuthEmail struct {
	repo repository.SendPassword
}

func NewSendPassword(repo repository.SendPassword) *AuthEmail {
	return &AuthEmail{repo: repo}
}

func GenerateTokenEmail() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *AuthEmail) CreateResetToken(email string) (string, error) {
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

	resetLink := fmt.Sprintf("http://localhost:8080/help/reset-password?token=%s", token)

	// Теперь отправляем email с ссылкой для сброса пароля
	subject := "Сброс пароля"
	body := fmt.Sprintf("Вы запросили сброс пароля. Пожалуйста, перейдите по следующей ссылке, чтобы сбросить ваш пароль: %s", resetLink)

	// Получаем случайный адрес отправителя
	from, password, err := s.getRandomSender()
	if err != nil {
		return "", err
	}

	if err := s.sendEmail(from, password, email, subject, body); err != nil {
		return "", err
	}

	return resetLink, nil
}

func (s *AuthEmail) getRandomSender() (string, string, error) {
	// Список адресов отправителей и их паролей
	senders := []struct {
		Email    string
		Password string
	}{
		{os.Getenv("EMAIL1"), os.Getenv("EMAIL1_PASSWORD")},
		{os.Getenv("EMAIL2"), os.Getenv("EMAIL2_PASSWORD")},
		{os.Getenv("EMAIL3"), os.Getenv("EMAIL3_PASSWORD")},
	}

	// Генерация случайного индекса
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(senders))

	if senders[index].Email == "" || senders[index].Password == "" {
		return "", "", fmt.Errorf("не удалось получить адрес отправителя или пароль")
	}

	return senders[index].Email, senders[index].Password, nil
}

func (s *AuthEmail) sendEmail(from string, password string, to string, subject string, body string) error {
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
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Printf("Ошибка при отправке email: %s", err)
		return err
	}

	log.Printf("Инструкция по восстановлению отправлена на почту %s", to)
	return nil
}
