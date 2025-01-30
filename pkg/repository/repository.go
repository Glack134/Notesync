package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/polyk005/notesync"
)

type Authorization interface {
	CreateUser(user notesync.User) (int, error)
	GetUser(username string, password string) (notesync.User, error)
	UpdatePasswordUser(username, newPasswordHash string) (notesync.User, error)
	GetUserIDByToken(token string) (int, error)
	UpdatePasswordUserByID(userID int, newPasswordHash string) error
	MarkTokenAsUsed(token string) error
	IsTokenUsed(token string) (bool, error)
	GetLastSentTime(token string) (time.Time, error)
}

type SendPassword interface {
	GetTokenResetPassword(email string) (int, error)
	SaveResetToken(userID int, token string, expiry time.Time) error
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
	GetById(userId, itemId int) (notesync.NotesyncItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input notesync.UpdateItemInput) error
}

type Repository struct {
	Authorization
	SendPassword
	NotesyncList
	NotesyncItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		SendPassword:  NewResetPostgres(db),
		NotesyncList:  NewNotesynsListPostgres(db),
		NotesyncItem:  NewNotesyncItemPostgres(db),
	}
}
