package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"my_notes_project/internal/entities"
)

type SQLiteDatabase struct {
	db *sql.DB
	/* TODO: logger */
}

func (d SQLiteDatabase) CloseSQLiteDatabase() error {
	return d.db.Close()
}

func (s SQLiteDatabase) AddUser(user *entities.User) (uint64, error) {
	res, err := s.db.Exec(`INSERT INTO User ("Name", "Password") VALUES ($1, $2);`, user.Name, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return uint64(id), err
}

func (s SQLiteDatabase) AddNote(note *entities.Note) (uint64, error) {
	res, err := s.db.Exec(`INSERT INTO Note ("Title", "Content", "UserID") VALUES ($1, $2, $3);`, note.Title, note.Content, note.UserID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return uint64(id), err
}

func (s SQLiteDatabase) RemoveNoteByID(ID uint64) error {
	_, err := s.db.Exec(`DELETE FROM Note WHERE id = $1;`, ID)

	return err
}

func (s SQLiteDatabase) UpdateNote(note *entities.Note) error {
	_, err := s.db.Exec(`UPDATE Note SET Title = $1, Content = $2 WHERE id = $3;`, note.Title, note.Content, note.ID)

	return err
}

func (s SQLiteDatabase) GetAllNotes() (map[uint64]*entities.Note, error) {
	row, err := s.db.Query(`SELECT * FROM Note;`)
	if err != nil {
		return nil, err
	}

	notes := map[uint64]*entities.Note{}

	for row.Next() {
		note := &entities.Note{}
		if err := row.Scan(&note.ID, &note.Title, &note.Content, &note.UserID); err != nil {
			return nil, err
		}
		notes[note.ID] = note
	}

	return notes, nil
}

func (s SQLiteDatabase) GetUserByName(name string) (*entities.User, error) {
	row := s.db.QueryRow(`SELECT * FROM User WHERE Name = $1;`, name)

	user := &entities.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Password); err != nil {
		return nil, err
	}
	return user, nil
}

func (s SQLiteDatabase) GetNotesByUserName(username string) (map[uint64]*entities.Note, error) {
	row, err := s.db.Query(`
	SELECT Note.*
	FROM User
	JOIN Note
	ON User.ID = Note.UserID AND User.Name = $1;`, username)
	if err != nil {
		return nil, err
	}

	notes := map[uint64]*entities.Note{}
	for row.Next() {
		n := &entities.Note{}
		if err := row.Scan(&n.ID, &n.Title, &n.Content, &n.UserID); err != nil {
			return nil, err
		}
		notes[n.ID] = n
	}

	return notes, nil
}
