package learn

import (
	"errors"

	"golang.org/x/net/context"
)

type UserService interface {
	CreateUser(cxt context.Context, user *User) (*User, error)
	GetUser(cxt context.Context, id string) (*User, error)
}

type basicService struct {
	Users map[string]*User
}

func (s basicService) CreateUser(_ context.Context, user *User) (*User, error) {
	s.Users[user.Id] = user

	return user, nil
}

func (s basicService) GetUser(_ context.Context, id string) (*User, error) {
	user, ok := s.Users[id]
	if !ok {
		return nil, errors.New("Could not find user")
	}

	return user, nil
}

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Username  string
}
