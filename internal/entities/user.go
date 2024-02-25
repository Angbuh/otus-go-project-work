package entities

type User struct {
	ID       uint64
	Name     string
	Password string
}

func NewUser(name, password string) *User {
	return &User{
		Name:     name,
		Password: password,
	}
}
