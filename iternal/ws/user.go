package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

type message struct {
	User    *User
	Content []byte
}

type User struct {
	ID      int
	Conn    *websocket.Conn
	Message chan *message
}

func (User *User) read(lb *Lobby) {
	defer func() {
		User.Conn.Close()
	}()
	for {
		select {
		case msg := <-User.Message:
			if err := User.Conn.WriteMessage(1, msg.Content); err != nil {
				log.Errorf("Error writing message: %v", err)
			}
		}
	}
}

func (User *User) write(lb *Lobby) {
	defer User.Conn.Close()

	for {
		_, data, err := User.Conn.ReadMessage()
		if err != nil {
			return
		}

		msg := &message{
			User:    User,
			Content: data,
		}

		lb.Broadcast <- msg
	}
}
