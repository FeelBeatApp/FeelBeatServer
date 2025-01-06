package room

import "github.com/feelbeatapp/feelbeatserver/internal/lib"

type Room struct {
	id       string
	ownerId  string
	settings RoomSettings
	songs    []lib.Song
}

func NewRoom(id string, ownerId string, settings RoomSettings, songs []lib.Song) Room {
	return Room{
		id:       id,
		settings: settings,
		ownerId:  ownerId,
		songs:    songs,
	}
}

func (r Room) Id() string {
	return r.id
}
