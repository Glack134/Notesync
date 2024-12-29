package service

import (
	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/repository"
)

type NotesyncItemService struct {
	repo     repository.NotesyncItem
	listRepo repository.NotesyncList
}

func NewNotesyncItemService(repo repository.NotesyncItem, listRepo repository.NotesyncList) *NotesyncItemService {
	return &NotesyncItemService{repo: repo, listRepo: listRepo}
}

func (s *NotesyncItemService) Create(userId, listId int, item notesync.NotesyncItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, nil
	}
	return s.repo.Create(listId, item)
}

func (s *NotesyncItemService) GetAll(userId, listId int) ([]notesync.NotesyncItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *NotesyncItemService) GetById(userId, itemId int) (notesync.NotesyncItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *NotesyncItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *NotesyncItemService) Update(userId, itemId int, input notesync.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
