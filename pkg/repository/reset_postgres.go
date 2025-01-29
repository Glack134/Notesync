package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polyk005/notesync"
)

type ResetPostgres struct {
	db *sqlx.DB
}

func NewResetPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *ResetPostgres) UpdatePasswordUser(username, newPasswordHash string) (notesync.User, error) {
	var updatepassword notesync.User
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE username=$2", usersTable)
	_, err := r.db.Exec(query, newPasswordHash, username)
	if err != nil {
		return updatepassword, err
	}
	return updatepassword, nil
}

// Реализация метода CreateResetToken
func (r *ResetPostgres) GetTokenResetPassword(email string) (int, error) {
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
func (r *ResetPostgres) SaveResetToken(userID int, token string, expiry time.Time) error {
	query := "INSERT INTO reset_tokens (user_id, token, expiry) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, userID, token, expiry)
	return err
}
