package core

import (
	"fmt"
	"my_notes_project/internal/entities"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type FakeDatabase struct {
	notes      map[uint64]*entities.Note
	users      map[uint64]*entities.User
	nextUserID *uint64
	nextNoteID *uint64
}

func NewFakeDatabase() *FakeDatabase {
	var uid uint64 = 0
	var nid uint64 = 0
	return &FakeDatabase{
		notes:      map[uint64]*entities.Note{},
		users:      map[uint64]*entities.User{},
		nextUserID: &uid,
		nextNoteID: &nid,
	}
}

func (f FakeDatabase) AddUser(user *entities.User) (uint64, error) {
	if user.Name == "" {
		return 0, fmt.Errorf("empty username")
	}

	user.ID = *f.nextUserID
	f.users[user.ID] = user
	*f.nextUserID += 1

	return 0, nil
}

func (f FakeDatabase) AddNote(note *entities.Note) (uint64, error) {
	if _, exists := f.users[note.UserID]; !exists {
		return 0, fmt.Errorf("user doesn't exist")
	}

	note.ID = *f.nextNoteID
	*f.nextNoteID += 1
	f.notes[note.ID] = note

	return 0, nil
}

func (f FakeDatabase) RemoveNoteByID(id uint64) error {
	delete(f.notes, id)

	return nil
}

func (f FakeDatabase) UpdateNote(note *entities.Note) error {
	f.notes[note.ID] = note

	return nil
}

func (f FakeDatabase) GetAllNotes() (map[uint64]*entities.Note, error) {

	return f.notes, nil
}

func (f FakeDatabase) GetUserByName(name string) (*entities.User, error) {
	for _, u := range f.users {
		if u.Name == name {
			return u, nil
		}
	}

	return nil, nil
}

func (f FakeDatabase) GetNotesByUserName(userName string) (map[uint64]*entities.Note, error) {
	var user *entities.User
	if u, err := f.GetUserByName(userName); u == nil || err != nil {
		return nil, err
	} else {
		user = u
	}

	notes := map[uint64]*entities.Note{}

	for _, note := range f.notes {
		if note.UserID == user.ID {
			notes[note.ID] = note
		}
	}

	return notes, nil

}

func TestRegisterUser(t *testing.T) {
	db := NewFakeDatabase()
	log := logrus.New()
	core := NewTheCore(db, log)

	expectedUser0 := entities.User{
		Name:     "Ivan",
		Password: "123",
	}

	assert.Equal(t, 0, len(db.users))
	assert.Equal(t, 0, len(db.notes))

	err := core.RegisterUser(expectedUser0.Name, expectedUser0.Password, expectedUser0.Password)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(db.users))
	assert.Equal(t, 0, len(db.notes))
	assert.Equal(t, *db.users[0], expectedUser0)

	expectedUser1 := entities.User{
		Name:     "Nikolay",
		Password: "321",
	}

	err = core.RegisterUser(expectedUser1.Name, expectedUser1.Password, expectedUser1.Password)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(db.users))
	assert.Equal(t, 0, len(db.notes))
	assert.Equal(t, *db.users[0], expectedUser0)
	assert.Equal(t, *db.users[1], expectedUser1)
}

func TestNotEqualPasswords(t *testing.T) {
	db := NewFakeDatabase()
	log := logrus.New()
	core := NewTheCore(db, log)

	err := core.RegisterUser("Ivan", "123", "321")

	assert.NotNil(t, err)
	assert.Equal(t, db, NewFakeDatabase())
}

func TestInvalidUserName(t *testing.T) {
	db := NewFakeDatabase()
	log := logrus.New()
	core := NewTheCore(db, log)

	err := core.RegisterUser("", "123", "123")

	assert.NotNil(t, err)
	assert.Equal(t, db, NewFakeDatabase())
}

func TestUpdateNoteByUserName(t *testing.T) {
	db := NewFakeDatabase()
	log := logrus.New()

	u := entities.User{
		ID:       3,
		Name:     "Ivan",
		Password: "123",
	}
	n := entities.Note{
		ID:      1,
		Title:   "Beach",
		Content: "nice beach and ocean",
		UserID:  u.ID,
	}

	expectedNote := entities.Note{
		ID:      n.ID,
		Title:   "Linux Modules",
		Content: "Linux is the best OS",
		UserID:  u.ID,
	}

	ns := map[uint64]*entities.Note{
		0: {
			ID:      0,
			Title:   "Tree",
			Content: "One,two",
			UserID:  u.ID,
		},
		n.ID: &n,
	}

	db.notes = ns
	db.users = map[uint64]*entities.User{
		u.ID: &u,
	}

	core := NewTheCore(db, log)

	assert.Equal(t, &n, db.notes[n.ID])
	err := core.UpdateNoteByUserName(u.Name, &expectedNote)
	assert.Nil(t, err)
	assert.Equal(t, &expectedNote, db.notes[expectedNote.ID])
	core.UpdateNoteByUserName(u.Name, &expectedNote)

}

func TestEmptyTitleAndContent(t *testing.T) {
	db := NewFakeDatabase()
	log := logrus.New()
	core := NewTheCore(db, log)

	u := entities.User{
		ID:       3,
		Name:     "Ivan",
		Password: "123",
	}

	n := entities.Note{
		ID:      1,
		Title:   "Beach",
		Content: "nice beach and ocean",
		UserID:  u.ID,
	}

	db.notes = map[uint64]*entities.Note{
		n.ID: &n,
	}
	db.users = map[uint64]*entities.User{
		u.ID: &u,
	}

	note := entities.Note{
		ID:      n.ID,
		Title:   n.Title,
		Content: "",
		UserID:  u.ID,
	}

	err := core.UpdateNoteByUserName(u.Name, &note)

	assert.NotNil(t, err)

	note = entities.Note{
		ID:      n.ID,
		Title:   "",
		Content: n.Content,
		UserID:  u.ID,
	}

	err = core.UpdateNoteByUserName(u.Name, &note)

	assert.NotNil(t, err)
}

func TestIsValidUserCredentials(t *testing.T) {
	db := NewFakeDatabase()
	log := logrus.New()

	u := entities.User{
		ID:       3,
		Name:     "Ivan",
		Password: "123",
	}

	db.users = map[uint64]*entities.User{
		u.ID: &u,
	}

	core := NewTheCore(db, log)
	isValid, err := core.IsValidUserCredentials(u.Name, u.Password)

	assert.Nil(t, err)
	assert.True(t, isValid)

}
