package services

import (
	"github.com/goombaio/namegenerator"
	"time"
)

type User struct {
	ID       int64
	Name     string
	Sender   chan []byte
	Receiver chan []byte
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

func (b *UserService) CreateUser(name string) User {
	if name == "" {
		name = b.NameGenerator.Generate()
	}
	usr := User{
		Name: name,
		// fixme (100)
		Sender:   make(chan []byte, 100),
		Receiver: make(chan []byte, 100),
	}
	if len(b.DB.Users) == 0 {
		usr.ID = 10000
	} else {
		usr.ID = b.DB.Users[len(b.DB.Users)-1].ID + 1
	}
	b.DB.Users = append(b.DB.Users, usr)
	return usr
}

func (b *UserService) GetUsers() []User {
	return b.DB.Users
}

func (b *UserService) GetUser(id int64) User {
	for _, user := range b.DB.Users {
		if user.ID == id {
			return user
		}
	}
	return User{}
}
