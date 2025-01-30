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
func (r *ResetPostgres) GetTokenResetPassword(email string) (int, time.Time, error) {
	var userID int
	var lastSent time.Time
	row := r.db.QueryRow("SELECT user_id, last_sent FROM users WHERE email = $1", email)
	err := row.Scan(&userID, &lastSent)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, time.Time{}, fmt.Errorf("пользователь с таким email не найден")
		}
		return 0, time.Time{}, err
	}
	return userID, lastSent, nil
}
