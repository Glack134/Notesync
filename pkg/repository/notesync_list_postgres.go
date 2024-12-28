package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/polyk005/notesync"
)

type NotesyncListPostgres struct {
	db *sqlx.DB
}

func NewNotesynsListPostgres(db *sqlx.DB) *NotesyncListPostgres {
	return &NotesyncListPostgres{db: db}
}

func (r *NotesyncListPostgres) Create(userId int, list notesync.NotesyncList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", notesyncListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		return 0, nil
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1 $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *NotesyncListPostgres) GetAll(userId int) ([]notesync.NotesyncList, error) {

	var lists []notesync.NotesyncList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tk.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ur.user_id = $1",
		notesyncListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *NotesyncListPostgres) GetById(userId, listId int) (notesync.NotesyncList, error) {
	var list notesync.NotesyncList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
						 INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		notesyncListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}
