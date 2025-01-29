package notesync

import "time"

type User struct {
	Id         int    `json:"-" db:"id"`
	Name       string `json:"name" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" binding:"required"`
	ResetToken string `json:"reset_token" db:"reset_token"`
}

type TokenRecord struct {
	UserID int       `json:"-" db:"id"`
	Token  string    `json:"token" db:"token"`
	Expiry time.Time `json:"expiry" db:"expiry"`
}
