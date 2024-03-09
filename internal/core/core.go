package core

import (
	"fmt"
	"my_notes_project/internal/database"
	"my_notes_project/internal/entities"

	"github.com/sirupsen/logrus"
)

type ServiceCore interface {
	GetAllNotes() (map[uint64]*entities.Note, error)
	RemoveNoteByID(uint64) error
	UpdateNoteByUserName(string, *entities.Note) error
	GetNotesByUserName(string) (map[uint64]*entities.Note, error)
	RegisterUser(string, string, string) error
	IsValidUserCredentials(string, string) (bool, error)
	AddNoteToUserByName(string, *entities.Note) error
}

type TheCore struct {
	db     database.DBRepository
	logger *logrus.Logger
}

func NewTheCore(db database.DBRepository, logger *logrus.Logger) *TheCore {
	return &TheCore{
		db:     db,
		logger: logger,
	}
}

func (c TheCore) GetAllNotes() (map[uint64]*entities.Note, error) {
	return c.db.GetAllNotes()
}

func (c TheCore) RemoveNoteByID(id uint64) error {
	return c.db.RemoveNoteByID(id)
}

func (c TheCore) UpdateNoteByUserName(username string, note *entities.Note) error {
	notes, err := c.db.GetNotesByUserName(username)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	note, exists := notes[note.ID]
	if !exists {
		c.logger.Error("not found")
		return err
	}

	return c.db.UpdateNote(note)
}

func (c TheCore) GetNotesByUserName(username string) (map[uint64]*entities.Note, error) {
	return c.db.GetNotesByUserName(username)
}

func (c TheCore) RegisterUser(name, password, repeatedPassword string) error {
	if password != repeatedPassword {
		return fmt.Errorf("passwords do not match")
	}

	user := &entities.User{
		Name:     name,
		Password: password,
	}

	id, err := c.db.AddUser(user)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	user.ID = id

	c.logger.Debug(user)
	return nil
}

func (c TheCore) IsValidUserCredentials(username, password string) (bool, error) {
	var err error
	var user *entities.User
	if user, err = c.db.GetUserByName(username); err != nil {
		c.logger.Error(err)
		return false, err
	}

	if user.Password != password {
		c.logger.Error("invalid data")
		return false, err
	}
	return false, nil
}

func (c TheCore) AddNoteToUserByName(username string, note *entities.Note) error {
	user, err := c.db.GetUserByName(username)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	note.UserID = user.ID
	noteID, err := c.db.AddNote(note)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	c.logger.Debug(noteID)
	return nil
}
