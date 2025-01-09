package roomrepository

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
)

type RoomRepository interface {
	CreateRoom(user auth.User, settings lib.RoomSettings) (string, error)
	GetAllRooms() []*room.Room
	GetRoom(id string) *room.Room
}
