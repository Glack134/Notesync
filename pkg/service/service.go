package service

import "github.com/polyk005/notesync/pkg/repository"

type Authoriation interface {
}

type NotesyncList interface {
}

type NotesyncItem interface {
}

type Service struct {
	Authoriation
	NotesyncList
	NotesyncItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
