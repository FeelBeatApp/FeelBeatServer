package room

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/messages"
)

type Room struct {
	id       string
	playlist lib.PlaylistData
	owner    lib.UserProfile
	settings RoomSettings
	players  map[string]Player
	hub      messages.Hub
	snd      chan messages.ServerMessage
	rcv      <-chan messages.ClientMessage
}

type Player struct {
	profile lib.UserProfile
}

func NewRoom(id string, playlist lib.PlaylistData, owner lib.UserProfile, settings RoomSettings, hub messages.Hub) *Room {
	return &Room{
		id:       id,
		playlist: playlist,
		owner:    owner,
		settings: settings,
		players:  make(map[string]Player),
		hub:      hub,
		snd:      make(chan messages.ServerMessage),
	}
}

func (r *Room) Id() string {
	return r.id
}

func (r *Room) Name() string {
	return r.playlist.Name
}

func (r *Room) PlayerProfiles() []lib.UserProfile {
	profiles := make([]lib.UserProfile, 0)
	for _, p := range r.players {
		profiles = append(profiles, p.profile)
	}

	return profiles
}

func (r *Room) ImageUrl() string {
	return r.playlist.ImageUrl
}

func (r *Room) Settings() RoomSettings {
	return r.settings
}

func (r *Room) Start() {
	r.rcv = r.hub.Run(r.snd)
	go r.processMessages()
}

func (r *Room) Stop() {
	close(r.snd)
}

func (r *Room) Hub() messages.Hub {
	return r.hub
}

func (r *Room) processMessages() {
	for message := range r.rcv {
		fblog.Info(component.Room, "Received message", "room", r.id, "from", message.From, "type", message.Type, "payload", message.Payload)
	}
}
