package controllers

import (
	"context"
	"encoding/json"
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
		defer conn.Close()

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

		err = greeting(conn, room, user)
		if err != nil {
			return
		}

		go transitOut(conn, user)
		go transitIn(conn, user)

		select {
		case <-ctx.Done():
			return
		}
	})
}

func transitIn(conn *websocket.Conn, user *models.User) {
	defer func() {
		conn.Close()
	}()

	for {
		response := new(lib.Request)
		err := conn.ReadJSON(response)
		if err != nil {
			log.Errorf("Cannot read response: %s", err)
			return
		}
		response.UID = user.Id()

		*user.Sender() <- *response
	}
}

func transitOut(conn *websocket.Conn, user *models.User) {
	defer func() {
		conn.Close()
	}()

	for {
		msg := <-*user.Receiver()
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Errorf("Cannot write response: %s", err)
			return
		}
	}
}

func greeting(conn *websocket.Conn, room *models.Room, user *models.User) error {
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

	pckHeaders, err := getPckHeaders(room)
	if err != nil {
		return err
	}
	err = conn.WriteJSON(pckHeaders)
	if err != nil {
		return err
	}

	return nil
}

func getPckHeaders(room *models.Room) (lib.Pack, error) {
	p := room.Pack()
	marshal, err := json.Marshal(p)
	if err != nil {
		return lib.Pack{}, err
	}

	siPckHeaders := new(lib.Pack)
	err = json.Unmarshal(marshal, &siPckHeaders)
	if err != nil {
		return lib.Pack{}, err
	}

	for _, round := range siPckHeaders.Rounds {
		for _, theme := range round.Themes {
			for _, question := range theme.Questions {
				question.Type = nil
				question.AnswerSlides = nil
				question.QuestionSlides = nil
				question.PriceMax = nil
				question.PriceStep = nil
			}
		}
	}
	return *siPckHeaders, nil
}
