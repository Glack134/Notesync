package repository

import "github.com/jmoiron/sqlx"

type Authoriation interface {
}

type NotesyncList interface {
}

type NotesyncItem interface {
}

type Repository struct {
	Authoriation
	NotesyncList
	NotesyncItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
