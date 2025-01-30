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
	// Получаем userID для email
	userID, err := s.repo.GetTokenResetPassword(email)
	if err != nil {
		log.Printf("Ошибка получения userID для email %s: %s", email, err)
		return "", err
	}

	// Генерируем токен
	token, err := GenerateTokenEmail()
	if err != nil {
		log.Printf("Ошибка генерации токена для email %s: %s", email, err)
		return "", err
	}

	// Устанавливаем срок действия токена
	expiry := time.Now().Add(tokenEmail)
	if err := s.repo.SaveResetToken(userID, token, expiry); err != nil {
		log.Printf("Ошибка сохранения токена для userID %d: %s", userID, err)
		return "", err
	}

	// Формируем ссылку для сброса пароля
	resetLink := fmt.Sprintf("http://localhost:8080/help/reset-password?token=%s", token)

	// Получаем случайный адрес отправителя
	from, password, err := s.getRandomSender()
	if err != nil {
		return "", err
	}

	// Подготовка HTML-сообщения
	body := fmt.Sprintf(`
		<html>
<head>
    <title>Сброс пароля Steam</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            border: 1px solid #eaeaea;
            border-radius: 5px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            overflow: hidden;
        }
        .header {
            background-color: #0d1228;
            color: #ffffff;
            padding: 20px;
            text-align: center;
        }
        .header img {
            width: 150px;
        }
        .content {
            padding: 20px;
            color: #333333;
        }
        .content h2 {
            color: #0d1228;
        }
        .reset-link {
            display: inline-block;
            margin: 20px 0;
            padding: 10px 20px;
            background-color:rgb(31, 90, 199);
            color: white;
            text-decoration: none;
            border-radius: 5px;
        }
        .footer {
            background-color: #f9f9f9;
            padding: 10px;
            text-align: center;
            font-size: 12px;
            color: #777777;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="https://i.pinimg.com/736x/ea/17/b4/ea17b4f210a18116b1d6dcca103a79c7.jpg" alt="Notesync_Logo"/>
        </div>
        <div class="content">
            <h2>Здравствуйте!</h2>
            <p>Вы запросили сброс пароля для вашей учетной записи Notesync. Пожалуйста, нажмите на кнопку ниже, чтобы сбросить ваш пароль.</p>
            <a href="%s" class="reset-link">Сбросить пароль</a>
            <p>Если вы не запрашивали сброс пароля, просто проигнорируйте это письмо.</p>
        </div>
        <div class="footer">
            <p>С уважением,<br>Команда Notesync</p>
        </div>
    </div>
</body>
</html>
	`, resetLink)

	// Отправляем email с ссылкой для сброса пароля
	subject := "Сброс пароля"
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
	message := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"To: " + to + "\r\n" +
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
