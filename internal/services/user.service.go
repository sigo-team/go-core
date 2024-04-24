package services

import (
	"github.com/goombaio/namegenerator"
	"sigo/internal/lib"
	"time"
)

type User struct {
	ID       int64
	Name     string
	Sender   chan lib.Request
	Receiver chan lib.Response
}

type UserService struct {
	NameGenerator     namegenerator.Generator
	IdentifierManager *lib.IdentifierManager
	DB                struct {
		Users map[int64]*User
	}
}

func NewUserService() *UserService {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	return &UserService{
		NameGenerator:     nameGenerator,
		IdentifierManager: lib.NewIdentifierManager(),
		DB:                struct{ Users map[int64]*User }{Users: make(map[int64]*User)},
	}
}

func (b *UserService) CreateUser(name string) *User {
	if name == "" {
		name = b.NameGenerator.Generate()
	}
	userId := b.IdentifierManager.NewID()
	usr := &User{
		Name: name,
		ID:   userId,
		// fixme (100)
		Sender:   make(chan lib.Request, 100),
		Receiver: make(chan lib.Response, 100),
	}

	b.DB.Users[userId] = usr
	return usr
}

func (b *UserService) GetUser(userId int64) *User {
	return b.DB.Users[userId]
}
