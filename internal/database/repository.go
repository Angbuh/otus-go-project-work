package database

import (
	"database/sql"
	"my_notes_project/internal/entities"
)

type DBRepository interface {
	AddUser(*entities.User) (uint64, error)
	AddNote(*entities.Note) (uint64, error)
	RemoveNoteByID(uint64) error
	UpdateNote(*entities.Note) error
	GetNoteById(ID uint64) (*entities.Note, error)
	GetAllNotes() ([]*entities.Note, error)
	GetUserByName(name string) (*entities.User, error)
	GetNotesByUserId(ID uint64) ([]*entities.Note, error)
	GetNotesByUserName(userName string) (map[uint64]*entities.Note, error)
}

func NewSQLiteDatabase(path string, /* TODO: logger */) (SQLiteDatabase, error) {
	d, err := sql.Open("sqlite3", "database/noteuser.db")
	if err != nil {
		return SQLiteDatabase{}, err
	}

	return SQLiteDatabase{
		db: d,
	}, nil
}
