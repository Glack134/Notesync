package repository

import (
	"database/sql"
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

func (r *AuthPostgres) FindUserByEmail(User int, email string) (notesync.User, error) {
	var user notesync.User
	query := "SELECT * FROM users WHERE id = $1 AND email = $2" // Предполагаем, что у вас есть колонка email

	err := r.db.QueryRow(query, user, email).Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil // Пользователь не найден
		}
		return user, err // Возвращаем ошибку, если произошла другая ошибка
	}

	return user, nil
}

func SendEmail(to string, subject string, body string) error {
	// Реализуйте логику отправки email здесь
	fmt.Printf("Sending email to: %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}
