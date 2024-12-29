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
	Delete(userId, listId int) error
	Update(userId, listId int, input notesync.UpdateListInput) error
}

type NotesyncItem interface {
	Create(listId int, item notesync.NotesyncItem) (int, error)
	GetAll(userId, listId int) ([]notesync.NotesyncItem, error)
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
		NotesyncItem:  NewNotesyncItemPostgres(db),
	}
}
