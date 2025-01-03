package roomrepository

import "github.com/feelbeatapp/feelbeatserver/internal/lib/room"

type RoomRepository interface {
	CreateRoom(string, room.RoomSettings) (string, error)
}
