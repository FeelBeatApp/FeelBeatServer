package roomapi

import (
	"github.com/feelbeatapp/feelbeatserver/internal/lib/roomrepository"
)

type RoomApi struct {
	roomRepo roomrepository.RoomRepository
}

func NewRoomApi(roomRepo roomrepository.RoomRepository) RoomApi {
	return RoomApi{
		roomRepo: roomRepo,
	}
}
