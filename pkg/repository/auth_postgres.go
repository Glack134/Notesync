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
	queryCheck := fmt.Sprintf("SELECT id FROM %s WHERE email=$1", usersTable)
	var userID int
	err := r.db.QueryRow(queryCheck, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если не найдено ни одной строки, возвращаем ошибку, что пользователь не существует
			return "", fmt.Errorf("user with email %s does not exist", email)
		}
		// Если произошла другая ошибка при выполнении запроса, возвращаем её
		return "", err
	}

	return email, nil // Возвращаем сгенерированный токен
}

// Реализация метода GetEmailByResetToken
func (r *AuthPostgres) GetEmailByResetToken(token string) (string, error) {
	var email string

	// Запрос к базе данных для получения email по токену
	query := "SELECT email FROM reset_tokens WHERE token = $1"
	err := r.db.QueryRow(query, token).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("token not found")
		}
		return "", err
	}

	return email, nil
}

// Реализация метода UpdatePassword
func (r *AuthPostgres) UpdatePassword(email, newPassword string) error {
	return nil
}
