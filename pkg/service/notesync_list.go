package service

import (
	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/repository"
)

type NotesyncListService struct {
	repo repository.NotesyncList
}

func NewNotesyncListService(repo repository.NotesyncList) *NotesyncListService {
	return &NotesyncListService{repo: repo}
}

func (s *NotesyncListService) Create(userId int, list notesync.NotesyncList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *NotesyncListService) GetAll(userId int) ([]notesync.NotesyncList, error) {
	return s.repo.GetAll(userId)
}

func (s *NotesyncListService) GetById(userId, listId int) (notesync.NotesyncList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *NotesyncListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *NotesyncListService) Update(userId, listId int, input notesync.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
