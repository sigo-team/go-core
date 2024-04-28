package controllers

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/utils"
	"os"
	"sigo/internal/lib"
	"sigo/internal/models"
	"sigo/internal/services"
	"strconv"
)

type RoomController struct {
	roomService *services.RoomService
}

type RoomControllerOptions struct {
	RoomService *services.RoomService
}

type UserController struct {
	userService *services.UserService
}

type UserControllerOptions struct {
	UserService *services.UserService
}

func validateRoomControllerOptions(options RoomControllerOptions) error {
	if options.RoomService == nil {
		return errors.New("")
	}
	return nil
}

func NewRoomController(options RoomControllerOptions) *RoomController {
	err := validateRoomControllerOptions(options)
	if err != nil {
		panic(err)
	}
	return &RoomController{
		roomService: options.RoomService,
	}
}

func validateUserControllerOptions(options UserControllerOptions) error {
	if options.UserService == nil {
		return errors.New("")
	}
	return nil
}

func NewUserController(options UserControllerOptions) *UserController {
	err := validateUserControllerOptions(options)
	if err != nil {
		panic(err)
	}
	return &UserController{
		userService: options.UserService,
	}
}

const UserIDKey = "user-id-key"

func (r *RoomController) CreateRoom(ctx *fiber.Ctx) error {
	var roomConfig models.RoomConfig
	err := ctx.BodyParser(&roomConfig)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	formFile, err := ctx.FormFile("file")
	if err != nil {
		log.Error(err)
		return err
	}
	packageName := utils.UUIDv4()
	err = ctx.SaveFile(formFile, "./tmp_"+packageName)
	if err != nil {
		log.Error(err)
		return err
	}
	err = lib.Unzip(packageName)
	if err != nil {
		log.Error(err)
		return err
	}
	err = os.Remove("./tmp_" + packageName)
	if err != nil {
		return err
	}

	user := ctx.Locals(UserIDKey).(*models.User)

	room, err := r.roomService.CreteRoom(models.RoomOptions{
		Config:      roomConfig,
		Owner:       user,
		PackageName: packageName,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(struct {
		RoomId int64 `json:"room_id"`
	}{room.Id()})
}

func (r *RoomController) GetRooms(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return err
	}
	rooms, err := r.roomService.ReadRooms(page)
	if err != nil && page != 0 {
		log.Error(err)
		return err
	}
	pages := r.roomService.GetRoomsAmount()/8 + 1
	return ctx.JSON(struct {
		Pages int            `json:"pages"`
		Rooms []*models.Room `json:"rooms"`
	}{pages, rooms})
}
