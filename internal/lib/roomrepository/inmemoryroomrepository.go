package roomrepository

import (
	"fmt"

	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/validation"
	"github.com/google/uuid"
)

type InMemoryRoomRepository struct {
	rooms map[string]room.Room
}

func NewInMemoryRoomRepository() InMemoryRoomRepository {
	return InMemoryRoomRepository{
		rooms: make(map[string]room.Room),
	}
}

// TODO: Implement fetching playlist details from spotify
func (r InMemoryRoomRepository) CreateRoom(ownderId string, settings room.RoomSettings) (string, error) {
	err := validation.ValidateRoomSettings(settings)
	if err != nil {
		return "", err
	}

	newRoom := room.NewRoom(uuid.NewString(), ownderId, settings)
	r.rooms[newRoom.Id()] = newRoom

	fmt.Println(r.rooms)

	return newRoom.Id(), nil
}
