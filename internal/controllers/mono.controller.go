package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"mime/multipart"
	"sigo/internal/services"
	"strconv"
)

type MonoService interface {
	CreateRoom(services.Room) (int64, error)
	GetRooms(int, string) (map[int64]*services.Room, int, error)
	//GetRoom(int64) (*services.Room, error)
}

type UserService interface {
	CreateUser(string) *services.User
	GetUser(int64) *services.User
}

type RoomHandlers struct {
	MonoService MonoService
	UserService UserService
}

func (r *RoomHandlers) CreateRoom(ctx *fiber.Ctx) error {
	var (
		err  error
		room services.Room
		file *multipart.FileHeader
	)

	err = ctx.BodyParser(&room)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	file, err = ctx.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = ctx.SaveFile(file, fmt.Sprintf("./%s", file.Filename))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	room.PackageName = file.Filename

	userId := ctx.Locals(UserIDKey).(int64)

	room.Khil = userId
	roomId, err := r.MonoService.CreateRoom(room)
	if err != nil {
		return err
	}

	ctx.Locals("roomId", roomId)
	if err := ctx.Next(); err != nil {
		log.Errorf("Cannot go ctx.next: %s", err.Error())
	}

	return ctx.SendString(fmt.Sprintf("Room created, id: %d", room.ID))
}

func (r *RoomHandlers) GetRooms(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return err
	}
	key := ctx.Query("key")
	rooms, pages, err := r.MonoService.GetRooms(page, key)
	if err != nil {
		return err
	}
	return ctx.JSON(struct {
		Pages int                      `json:"pages"`
		Rooms map[int64]*services.Room `json:"rooms"`
	}{pages, rooms})
}
