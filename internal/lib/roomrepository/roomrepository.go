package roomrepository

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
)

type RoomRepository interface {
	CreateRoom(user auth.User, settings room.RoomSettings) (string, error)
	GetAllRooms() []room.Room
}
