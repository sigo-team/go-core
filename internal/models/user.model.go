package models

import (
	"sigo/internal/lib"
)

type User struct {
	id       int64
	name     string
	sender   *chan lib.Request
	receiver *chan lib.Response
}

func (user *User) Name() string {
	if user.name == "" {
		panic("using user name before mounting it to the service")
	}
	return user.name
}

func (user *User) Id() int64 {
	if user.id == 0 {
		panic("using user id before mounting it to the service")
	}
	return user.id
}

func (user *User) Mount(id int64, name string) {
	if user.id != 0 {
		panic("user is already mounted")
	}
	user.id = id
	user.name = name
}

func NewUser() (*User, error) {
	sender := make(chan lib.Request, 100)
	receiver := make(chan lib.Response, 100)
	return &User{
		// FIXME: (100)
		sender:   &sender,
		receiver: &receiver,
	}, nil
}

func (user *User) Sender() *chan lib.Request {
	return user.sender
}

func (user *User) SetSender(ch *chan lib.Request) {
	user.sender = ch
}

func (user *User) Receiver() *chan lib.Response {
	return user.receiver
}
