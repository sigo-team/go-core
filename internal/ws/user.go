package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/controllers"
	"sigo/internal/lib"
)

type message struct {
	UserID  string
	Content []byte
}

type User struct {
	ID       string
	Conn     *websocket.Conn
	Receiver chan *message
	Sender   chan *message
}

func (User *User) read(lb *Lobby) {
	defer func() {
		User.Conn.Close()
	}()
	for {
		select {
		case msg := <-User.Receiver:
			if err := User.Conn.WriteMessage(1, msg.Content); err != nil {
				log.Errorf("Error writing message: %v", err)
			}
		}
	}
}

func (user *User) write(lb *Lobby) error {
	defer func() {
		user.Conn.Close()
	}()
	for {
		_, data, err := user.Conn.ReadMessage()
		if err != nil {
			return err
		}

		response := controllers.ReadMessage(data)
		if response.Type == "pressButton" {
			log.Debug("Button pressed")
			msg := &message{
				UserID: user.ID,
				Content: lib.Request{
					Type: "pressButton",
					Data: lib.Data{
						PlayerId: user.ID,
					},
				}.Marshall(),
			}
			lb.ButtonBC <- msg
			return nil
		}

		msg := &message{
			UserID:  user.ID,
			Content: data,
		}
		user.Sender <- msg
	}
}
