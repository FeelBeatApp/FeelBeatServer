package room

import (
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
)

type Room struct {
	id       string
	playlist lib.PlaylistData
	owner    lib.UserProfile
	settings RoomSettings
	players  []Player
}

type Player struct {
}

func NewRoom(id string, playlist lib.PlaylistData, owner lib.UserProfile, settings RoomSettings) Room {
	return Room{
		id:       id,
		playlist: playlist,
		owner:    owner,
		settings: settings,
		players:  make([]Player, 0),
	}
}

func (r Room) Id() string {
	return r.id
}

func (r Room) Name() string {
	return r.playlist.Name
}

func (r Room) Players() []Player {
	return r.players
}

func (r Room) ImageUrl() string {
	return r.playlist.ImageUrl
}
func (r Room) Settings() RoomSettings {
	return r.settings
}
