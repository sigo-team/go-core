package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
	"sigo/internal/lib"
	"sigo/internal/models"
	"sigo/internal/services"
	"strconv"
)

// ctx *fiber.Ctx
func Handler(ctx context.Context, rc *RoomController) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		//todo: user disconnect
		roomId, err := strconv.ParseInt(conn.Query("room_id"), 10, 64)
		if err != nil {
			log.Errorf("Cannot parse room id from query: %s", conn.Query("room_id"))
			conn.Close()
			return
		}

		user := conn.Locals(UserIDKey).(*models.User)
		room, err := rc.roomService.ReadRoom(roomId)
		if err != nil {
			log.Errorf("Cannot read room %d: %s", roomId, err)
			return
		}

		if user.Id() == room.Owner().Id() {
			log.Infof("Owner %d joined room %d", user.Id(), roomId)
			go services.Listening(room, ctx)
		} else {
			room.JoinPlayer(user)
			log.Infof("New player %d joined room %d", user.Id(), roomId)
		}

		conn.WriteJSON(room)

		go transitIn(conn, user)
		go transitOut(conn, user)

		select {
		case <-ctx.Done():
			conn.Close()
			return
		}
	})
}

func transitIn(conn *websocket.Conn, user *models.User) {
	for {
		response := new(lib.Request)
		err := conn.ReadJSON(response)
		if err != nil {
			log.Errorf("Cannot read response: %s", err)
			return
		}
		response.UID = user.Id()

		log.Debugf("Read response: %s", response)

		user.Sender() <- *response
	}
}

func transitOut(conn *websocket.Conn, user *models.User) {
	for {
		msg := <-user.Receiver()

		err := conn.WriteJSON(msg)
		if err != nil {
			log.Errorf("Cannot write response: %s", err)
			return
		}
	}
}
