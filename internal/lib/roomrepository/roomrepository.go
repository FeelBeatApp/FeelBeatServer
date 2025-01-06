package roomrepository

import (
	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
)

type RoomRepository interface {
	CreateRoom(playlistId string, settings room.RoomSettings, token string) (string, error)
}
