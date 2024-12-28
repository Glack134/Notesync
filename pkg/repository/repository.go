package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/polyk005/notesync"
)

type Authorization interface {
	CreateUser(user notesync.User) (int, error)
	GetUser(username string, password string) (notesync.User, error)
}

type NotesyncList interface {
	Create(userId int, list notesync.NotesyncList) (int, error)
	GetAll(userId int) ([]notesync.NotesyncList, error)
	GetById(userId, listId int) (notesync.NotesyncList, error)
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
		NotesyncList:  NewNotesynsListPostgres(db),
	}
}
