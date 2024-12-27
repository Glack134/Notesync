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
}

type NotesyncItem interface {
}

type Service struct {
	Authorization
	NotesyncList
	NotesyncItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
