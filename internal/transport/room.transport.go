package transport

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"sigo/internal/services"
)

// todo: create room
func ConnectKhil(us *services.UserService) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		defer conn.Close()

		khil := services.NewKhil(us)

		go transitIn(conn, &khil.User)
		go transitOut(conn, &khil.User)

		//todo: register, closing
	})
}
