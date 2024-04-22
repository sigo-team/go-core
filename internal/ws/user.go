package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/api"
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

		response := api.ReadResponse(data)
		if response.Type == "pressButton" {
			log.Debug("Button pressed")
			msg := &message{
				UserID: user.ID,
				Content: api.Request{
					Type: "pressButton",
					Data: api.Data{
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
