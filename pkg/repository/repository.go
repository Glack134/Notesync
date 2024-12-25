package repository

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

func NewRepository() *Repository {
	return &Repository{}
}
