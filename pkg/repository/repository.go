package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/polyk005/notesync"
)

type Authorization interface {
	CreateUser(user notesync.User) (int, error)
}

type NotesyncList interface {
}

type NotesyncItem interface {
}

type Repository struct {
	Authorization
	NotesyncList
	NotesyncItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
