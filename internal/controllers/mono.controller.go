package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"sigo/internal/services"
	"strconv"
)

type MonoService interface {
	CreateRoom(services.Room) (int64, error)
	GetRooms(int, string) ([]services.Room, int, error)
}

type RoomHandlers struct {
	Service MonoService
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
	_, err = r.Service.CreateRoom(room)
	if err != nil {
		return err
	}
	return ctx.SendString(fmt.Sprintf("Room created, id: %d", room.ID))
}

func (r *RoomHandlers) GetRooms(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return err
	}
	key := ctx.Query("key")
	rooms, pages, err := r.Service.GetRooms(page, key)
	if err != nil {
		return err
	}
	return ctx.JSON(struct {
		Pages int             `json:"pages"`
		Rooms []services.Room `json:"rooms"`
	}{pages, rooms})
}
