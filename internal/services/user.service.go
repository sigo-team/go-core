package services

import (
	"errors"
	"github.com/goombaio/namegenerator"
	"sigo/internal/lib"
	"sigo/internal/models"
)

type UserService struct {
	nameGenerator     namegenerator.Generator
	identifierManager *lib.IdentifierManager
	users             map[int64]*models.User
}

type UserServiceOptions struct {
	NameGenerator     namegenerator.Generator
	IdentifierManager *lib.IdentifierManager
}

func validateUserServiceOptions(options UserServiceOptions) error {
	return nil
}

func NewUserService(options UserServiceOptions) *UserService {
	err := validateUserServiceOptions(options)
	if err != nil {
		panic(err)
	}
	return &UserService{
		nameGenerator:     options.NameGenerator,
		identifierManager: options.IdentifierManager,
		users:             make(map[int64]*models.User),
	}
}

func (b *UserService) CreateUser() (*models.User, error) {
	user, err := models.NewUser()
	if err != nil {
		return nil, err
	}
	user.Mount(b.identifierManager.NewID(), b.nameGenerator.Generate())
	b.users[user.Id()] = user
	return user, nil
}

func (b *UserService) ReadUser(userId int64) (*models.User, error) {
	user, ok := b.users[userId]
	if !ok {
		return nil, errors.New("")
	}
	return user, nil
}

func (b *UserService) DeleteUser(userId int64) error {
	_, ok := b.users[userId]
	if !ok {
		return errors.New("")
	}
	delete(b.users, userId)
	return nil
}
