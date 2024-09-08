package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"sigo/internal/lib"
	"sigo/internal/models"
	"strconv"
)

// ctx *fiber.Ctx
func Handler(ctx context.Context, rc *RoomController) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		roomId, err := strconv.ParseInt(conn.Query("room_id"), 10, 64)
		if err != nil {
			log.Errorf("Cannot parse room id from query: %s", conn.Query("room_id"))
			conn.Close()
		}

		user := conn.Locals(UserIDKey).(*models.User)
		room, err := rc.roomService.ReadRoom(roomId)
		if err != nil {
			log.Errorf("Cannot read room %d: %s", roomId, err)
			return
		}

		defer func() {
			*user.Sender() <- lib.Request{
				UID:  user.Id(),
				Type: "disconnected",
			}
			room.DisjoinPlayer(user)
			conn.Close()
			log.Infof("Connection closed %s", user.Id())
		}()

		user.SetSender(room.Receiver())
		room.JoinPlayer(user)

		log.Infof("Player %d joined room %d", user.Id(), roomId)
		sendGreeting(conn, room, user)
		*user.Sender() <- lib.Request{
			UID:  user.Id(),
			Type: "connected",
		}

		go write(conn, user)
		read(conn, user)
	})
}

func read(conn *websocket.Conn, user *models.User) {
	for {
		response := new(lib.Request)
		err := conn.ReadJSON(response)
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
			return
		} else if err != nil {
			log.Errorf("Cannot read json: %s", err)
			continue
		}

		response.UID = user.Id()

		*user.Sender() <- *response
	}
}

func write(conn *websocket.Conn, user *models.User) {
	for {
		response := <-*user.Receiver()
		err := conn.WriteJSON(response)
		if err != nil {
			log.Errorf("Cannot write response: %s", err)
			return
		}
	}
}

func sendGreeting(conn *websocket.Conn, room *models.Room, user *models.User) error {
	err := conn.WriteJSON(room)
	if err != nil {
		return err
	}

	if user.Id() == room.Owner().Id() {
		err := conn.WriteJSON(room.Pack())
		if err != nil {
			return err
		}
		return nil
	}

	pckHeaders, err := lib.GetPckHeaders(room.Pack())
	if err != nil {
		return err
	}
	err = conn.WriteJSON(pckHeaders)
	if err != nil {
		return err
	}

	return nil
}
