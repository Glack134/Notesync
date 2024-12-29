package service

import (
	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/repository"
)

type Authorization interface {
	CreateUser(user notesync.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type NotesyncList interface {
	Create(userId int, list notesync.NotesyncList) (int, error)
	GetAll(userId int) ([]notesync.NotesyncList, error)
	GetById(userId, listId int) (notesync.NotesyncList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input notesync.UpdateListInput) error
}

type NotesyncItem interface {
	Create(userId int, listId int, item notesync.NotesyncItem) (int, error)
	GetAll(userId, listId int) ([]notesync.NotesyncItem, error)
}

type Service struct {
	Authorization
	NotesyncList
	NotesyncItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		NotesyncList:  NewNotesyncListService(repos.NotesyncList),
		NotesyncItem:  NewNotesyncItemService(repos.NotesyncItem, repos.NotesyncList),
	}
}
