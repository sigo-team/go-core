package transport

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/controllers"
	"sigo/internal/services"
)

func read(conn *websocket.Conn, user *services.User) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Error while reading from websocket: %s", err)
			continue
		}

		controllers.SendMessage(user, msg)
	}
}

func write(conn *websocket.Conn, user *services.User) {
	for {
		msg := controllers.ReadMessage(user)
		err := conn.WriteMessage(1, msg)
		if err != nil {
			log.Errorf("Error while writing to websocket: %s", err)
			continue
		}
	}
}
