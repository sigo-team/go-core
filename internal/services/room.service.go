package services

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"sigo/internal/lib"
	"sigo/internal/models"
	"sort"
)

type RoomService struct {
	identifierManager *lib.IdentifierManager
	rooms             map[int64]*models.Room
}

type RoomServiceOptions struct {
	IdentifierManager *lib.IdentifierManager
}

func validateRoomServiceOptions(options RoomServiceOptions) error {
	return nil
}

func NewRoomService(options RoomServiceOptions) *RoomService {
	err := validateRoomServiceOptions(options)
	if err != nil {
		panic(err)
	}
	return &RoomService{identifierManager: options.IdentifierManager, rooms: make(map[int64]*models.Room)}
}

func (r *RoomService) CreteRoom(options models.RoomOptions) (*models.Room, error) {
	room, err := models.NewRoom(options)
	if err != nil {
		return nil, err
	}
	room.Mount(r.identifierManager.NewID())
	r.rooms[room.Id()] = room
	return room, nil
}

func (r *RoomService) ReadRoom(roomId int64) (*models.Room, error) {
	room, ok := r.rooms[roomId]
	if !ok {
		return nil, errors.New("")
	}
	return room, nil
}

func partitioned(p int, n int, slice []int64) []int64 {
	if n < 0 {
		return []int64{}
	}
	start := (n - 1) * p
	end := n * p
	return slice[start:min(end, len(slice))]
}

func (r *RoomService) GetRoomsAmount() int {
	return len(r.rooms)
}

func (r *RoomService) ReadRooms(page int) ([]*models.Room, error) {
	ids := make([]int64, 0)
	for _, room := range r.rooms {
		ids = append(ids, room.Id())
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	part := partitioned(8, page, ids)
	rooms := make([]*models.Room, 0)
	for _, id := range part {
		rooms = append(rooms, r.rooms[id])
	}
	return rooms, nil
}

func Listening(r *models.Room, closingCtx context.Context) {
	log.Info("Start listnening room")
	for {
		select {
		case msg := <-r.Owner().Sender():
			{
				log.Debugf("got msg from owner: %s", msg)
				response := lib.Response{
					UID:  msg.UID,
					Type: msg.Type,
					Data: msg.Data,
				}
				for _, user := range r.Players() {
					user.Receiver() <- response
				}
			}

		case <-closingCtx.Done():
			return
		}
	}
}
