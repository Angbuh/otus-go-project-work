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
	GetAllNotes() ([]*entities.Note, error)
	GetUserByName(name string) (*entities.User, error)
	GetNotesByUserName(userName string) (map[uint64]*entities.Note, error)
}

func NewSQLiteDatabase(path string /* TODO: logger */) (*SQLiteDatabase, error) {
	d, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &SQLiteDatabase{
		db: d,
	}, nil
}
