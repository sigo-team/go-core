package transport

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/controllers"
	"sigo/internal/lib"
	"sigo/internal/services"
)

func transitIn(conn *websocket.Conn, user *services.User) {
	for {
		response := new(lib.Response)
		err := conn.ReadJSON(response)
		if err != nil {
			log.Errorf("Error while reading from websocket: %s", err)
			continue
		}

		controllers.SendMessage(conn, user, response)
	}
}

func transitOut(conn *websocket.Conn, user *services.User) {
	for {
		request := controllers.ReadMessage(user)
		err := conn.WriteJSON(request)
		if err != nil {
			log.Errorf("Error while writing to websocket: %s", err)
		}
	}
}
