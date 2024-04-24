package controllers

import (
	"github.com/gofiber/contrib/websocket"
	"sigo/internal/lib"
	"sigo/internal/services"
)

const UserIDKey = "user-id-key"

func SendMessage(conn *websocket.Conn, user *services.User, response *lib.Response) {
	uid := new(int64)
	conn.Locals(UserIDKey, uid)

	response.UID = *uid

	user.Receiver <- *response
}

func ReadMessage(user *services.User) lib.Request {
	return <-user.Sender
}
