package repository

import (
	"database/sql"
	"errors"
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
	query := fmt.Sprintf(`INSERT INTO %s (name, username, password_hash, email) VALUES ($1, $2, $3, $4) RETURNING id`, usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password, user.Email)
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

// Реализация метода CreateResetToken
func (r *AuthPostgres) CreateResetToken(email string) (string, error) {
	// Генерация токена (например, с использованием UUID или другого метода)
	token := generateToken() // Предполагается, что у вас есть функция для генерации токена

	// Сохранение токена в базу данных
	query := "INSERT INTO reset_tokens (email, token) VALUES ($1, $2)"
	_, err := r.db.Exec(query, email, token)
	if err != nil {
		return "", err // Возвращаем ошибку, если не удалось сохранить токен
	}

	return token, nil // Возвращаем сгенерированный токен
}

// Реализация метода GetEmailByResetToken
func (r *AuthPostgres) GetEmailByResetToken(token string) (string, error) {
	var email string

	// Запрос к базе данных для получения email по токену
	query := "SELECT email FROM reset_tokens WHERE token = $1"
	err := r.db.QueryRow(query, token).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("token not found") // Возвращаем ошибку, если токен не найден
		}
		return "", err // Возвращаем другую ошибку
	}

	return email, nil // Возвращаем найденный email
}

// Реализация метода UpdatePassword
func (r *AuthPostgres) UpdatePassword(email, newPassword string) error {
	// Обновление пароля в базе данных
	query := "UPDATE users SET password = $1 WHERE email = $2"
	_, err := r.db.Exec(query, newPassword, email)
	if err != nil {
		return err // Возвращаем ошибку, если не удалось обновить пароль
	}

	return nil // Возвращаем nil, если все прошло успешно
}

// Вспомогательная функция для генерации токена
func generateToken() string {
	// Ваша реализация генерации токена (например, с использованием UUID)
	return "some-generated-token" // Замените на реальную генерацию
}
