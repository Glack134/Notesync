package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polyk005/notesync"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user notesync.User) (int, error) {
	var id int
	// Оберните имя таблицы в двойные кавычки, чтобы избежать конфликтов с зарезервированными словами
	query := fmt.Sprintf(`INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (notesync.User, error) {
	var user notesync.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}

func (r *AuthPostgres) SendPasswordResetEmail(username string) (notesync.User, error) {
	var user notesync.User
	query := fmt.Sprintf("SELECT id, email FROM %s WHERE username=$1", usersTable) // Предполагается, что email хранится в базе данных
	err := r.db.Get(&user, query, username)
	if err != nil {
		return notesync.User{}, fmt.Errorf("пользователь не найден: %w", err) // Возвращаем ошибку, если пользователь не найден
	}

	// Генерация токена сброса пароля
	resetToken := generatePasswordResetToken(user.Id) // Реализуйте эту функцию для генерации безопасного токена
	resetLink := fmt.Sprintf("http://yourdomain.com/reset-password?token=%s", resetToken)

	// Отправка электронного письма для сброса пароля
	err = sendEmail(user.Email, "Сброс пароля", resetLink) // Предполагается, что user.Email доступен
	if err != nil {
		return notesync.User{}, fmt.Errorf("не удалось отправить электронное письмо: %w", err) // Возвращаем ошибку, если отправка письма не удалась
	}

	return user, nil // Возвращаем пользователя и nil в качестве ошибки, если все прошло успешно
}

func generatePasswordResetToken(userId int) string {
	// Реализуйте логику генерации токена здесь (например, используя JWT или случайную строку)
	return fmt.Sprintf("token-for-user-%d", userId) // Заглушка для реализации
}

// Реализуйте функцию sendEmail
func sendEmail(to string, subject string, body string) error {
	// Реализуйте вашу логику отправки электронной почты здесь
	return nil // Заглушка для реализации
}
