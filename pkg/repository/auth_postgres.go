package repository

import (
	"database/sql"
	"fmt"
	"time"

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

func (r *AuthPostgres) UpdatePasswordUser(username, newPasswordHash string) (notesync.User, error) {
	var updatepassword notesync.User
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE username=$2", usersTable)
	_, err := r.db.Exec(query, newPasswordHash, username)
	if err != nil {
		return updatepassword, err
	}
	return updatepassword, nil
}

// Реализация метода CreateResetToken
func (r *AuthPostgres) GetTokenResetPassword(email string) (int, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1", usersTable)
	var userID int
	err := r.db.QueryRow(query, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user with email %s does not exist", email)
		}
		return 0, err
	}
	return userID, nil
}

// Исправленный метод SaveResetToken
func (r *AuthPostgres) SaveResetToken(userID int, token string, expiry time.Time) error {
	query := "INSERT INTO reset_tokens (user_id, token, expiry) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, userID, token, expiry)
	return err
}

// Получаем userID по токену
func (r *AuthPostgres) GetUserIDByToken(token string) (int, error) {
	var userID int
	query := "SELECT user_id FROM reset_tokens WHERE token=$1 AND expiry > NOW()"
	err := r.db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return 0, err // Возвращаем ошибку, если токен недействителен
	}
	return userID, nil // Возвращаем userID
}

// Обновляем пароль пользователя по userID
func (r *AuthPostgres) UpdatePasswordUserByID(userID int, newPasswordHash string) error {
	query := "UPDATE users SET password_hash=$1 WHERE id=$2"
	_, err := r.db.Exec(query, newPasswordHash, userID)
	return err // Возвращаем ошибку, если обновление не удалось
}
