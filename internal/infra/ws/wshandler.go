package ws

import "github.com/feelbeatapp/feelbeatserver/internal/lib/roomrepository"

type WSHandler struct {
	roomRepo roomrepository.RoomRepository
}

func NewWSHandler(basePath string, roomRepo roomrepository.RoomRepository) WSHandler {
	return WSHandler{
		roomRepo: roomRepo,
	}
}
