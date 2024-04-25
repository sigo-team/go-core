package transport

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
	"sigo/internal/controllers"
	"strconv"
)

func PublicRoutes(app *fiber.App, roomController *controllers.RoomController) {
	route := app.Group("/api/v1")

	route.Get("/room", roomController.GetRooms)
	route.Post("/room", roomController.CreateRoom)

	route.Use("/ws", UpgradeMiddleware)

	route.Get("/ws", websocket.New(func(c *websocket.Conn) {
		roomId, err := strconv.ParseInt(c.Query("room_id"), 10, 64)
		if err != nil {
			c.Close()
			return
		}

		log.Println(roomId)

		var (
			mt  int
			msg []byte
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))
}
