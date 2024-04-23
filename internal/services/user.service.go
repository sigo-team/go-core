package services

import (
	"errors"
	"github.com/goombaio/namegenerator"
	"time"
)

type User struct {
	ID   int
	Name string
}

type UserService struct {
	NameGenerator namegenerator.Generator
	DB            struct {
		Users []User
	}
}

func NewUserService() *UserService {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	return &UserService{
		NameGenerator: nameGenerator,
		DB: struct {
			Users []User
		}{
			Users: make([]User, 0),
		},
	}
}

func (b *UserService) CreateUser(name string) (User, error) {
	if name == "" {
		name = b.NameGenerator.Generate()
	}
	usr := User{
		Name: name,
	}
	if len(b.DB.Users) == 0 {
		usr.ID = 10000
	} else {
		usr.ID = b.DB.Users[len(b.DB.Users)-1].ID + 1
	}
	b.DB.Users = append(b.DB.Users, usr)
	return usr, nil
}

func (b *UserService) GetUsers() ([]User, error) {
	return b.DB.Users, nil
}

var (
	NotFoundErr = errors.New("user not found")
)

func (b *UserService) GetUser(id int) (User, error) {
	for _, user := range b.DB.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, NotFoundErr
}
