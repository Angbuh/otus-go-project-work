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

// Возвращает все заметки
func (c TheCore) GetAllNotes() (map[uint64]*entities.Note, error) {
	//Обращаемся в базу данных и возвращаем все заметки
	return c.db.GetAllNotes()
}

func (c TheCore) RemoveNoteByID(id uint64) error {
	//Обращаемся в базу данных и удаляем заметку по id
	return c.db.RemoveNoteByID(id)
}

func (c TheCore) UpdateNoteByUserName(username string, note *entities.Note) error {
	// Получаем заметки по имени пользователя
	notes, err := c.db.GetNotesByUserName(username)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	// Проверяем на существование заметку по id
	_, exists := notes[note.ID]
	if !exists {
		c.logger.Error("not found")
		return fmt.Errorf("not found")
	}

	// проверка на пустые Title и Content
	if note.Title == "" || note.Content == "" {
		return fmt.Errorf("empty title or content")
	}

	// Обновляем заметку, если она существует
	return c.db.UpdateNote(note)
}

func (c TheCore) GetNotesByUserName(username string) (map[uint64]*entities.Note, error) {
	// Обращаемся в базу данных и получаем заметки по имени пользователя
	return c.db.GetNotesByUserName(username)
}

func (c TheCore) RegisterUser(name, password, repeatedPassword string) error {
	//Проверяем совпадение паролей
	if password != repeatedPassword {
		return fmt.Errorf("passwords do not match")
	}

	//Создаем экземпляр пользователя
	user := &entities.User{
		Name:     name,
		Password: password,
	}

	//Добовляем пользователя в базу данных и получаем его id
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
	//Получаем пользователя по имени
	user, err := c.db.GetUserByName(username)
	if err != nil {
		c.logger.Error(err)
		return false, err
	}

	// err == nil

	//Проверяем действительный ли пароль пользователя и тот, который он вводит
	return user.Password == password, nil
}

func (c TheCore) AddNoteToUserByName(username string, note *entities.Note) error {
	//Получаем пользователя из базы данных по его имени
	user, err := c.db.GetUserByName(username)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	//Устанавливаем у заметки id пользователя-автора
	note.UserID = user.ID

	//Добавляем заметку в базу данных
	noteID, err := c.db.AddNote(note)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	c.logger.Debug(noteID)
	return nil
}
