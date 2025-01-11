package room

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/messages"
)

type Room struct {
	id        string
	playlist  lib.PlaylistData
	owner     lib.UserProfile
	settings  lib.RoomSettings
	players   map[string]Player
	hub       messages.Hub
	snd       chan messages.ServerMessage
	rcv       <-chan messages.ClientMessage
	onCleanup func(*Room)
}

type Player struct {
	profile lib.UserProfile
}

func NewRoom(id string, playlist lib.PlaylistData, owner lib.UserProfile, settings lib.RoomSettings, hub messages.Hub, onCleanup func(*Room)) *Room {
	return &Room{
		id:        id,
		playlist:  playlist,
		owner:     owner,
		settings:  settings,
		players:   make(map[string]Player),
		hub:       hub,
		snd:       make(chan messages.ServerMessage),
		onCleanup: onCleanup,
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

func (r *Room) Settings() lib.RoomSettings {
	return r.settings
}

func (r *Room) Start() {
	r.rcv = r.hub.Run(r.snd)
	go r.processMessages()
}

func (r *Room) Hub() messages.Hub {
	return r.hub
}

func (r *Room) addPlayer(profile lib.UserProfile) {
	r.players[profile.Id] = Player{
		profile: profile,
	}

	fblog.Info(component.Room, "new player", "roomId", r.id, "userId", profile.Id)

	playerProfiles := make([]lib.UserProfile, 0)
	for _, p := range r.players {
		playerProfiles = append(playerProfiles, p.profile)
	}

	packedSongs := make([]messages.SongState, 0)
	for _, s := range r.playlist.Songs {
		packedSongs = append(packedSongs, messages.SongState{
			Id:          s.Id,
			Title:       s.Details.Title,
			Artist:      s.Details.Artist,
			ImageUrl:    s.Details.ImageUrl,
			DurationSec: int(s.Details.Duration.Seconds()),
		})
	}

	r.snd <- messages.ServerMessage{
		To:   []string{profile.Id},
		Type: messages.InitialMessage,
		Payload: messages.InitialGameState{
			Id:    r.id,
			Me:    profile.Id,
			Admin: r.owner.Id,
			Playlist: messages.PlaylistState{
				Name:     r.Name(),
				ImageUrl: r.ImageUrl(),
				Songs:    packedSongs,
			},
			Players:  playerProfiles,
			Settings: r.settings,
		},
	}
	r.sendToAllExcept(profile.Id, messages.NewPlayer, profile)
}

func (r *Room) removePlayer(id string) {
	if _, ok := r.players[id]; !ok {
		return
	}

	delete(r.players, id)
	recipents := make([]string, 0)
	for _, p := range r.players {
		recipents = append(recipents, p.profile.Id)
	}
	if len(r.players) == 0 {
		r.onCleanup(r)
		r.cleanup()
		return
	}

	if id == r.owner.Id {
		r.owner = r.players[recipents[0]].profile
		fblog.Info(component.Room, "admin transfered", "roomId", r.id, "from", id, "to", r.owner.Id)
	}

	fblog.Info(component.Room, "player leaves", "roomId", r.id, "playerId", id)

	r.snd <- messages.ServerMessage{
		To:   recipents,
		Type: messages.PlayerLeft,
		Payload: messages.PlayerLeftPayload{
			Left:  id,
			Admin: r.owner.Id,
		},
	}
}

func (r *Room) sendToAllExcept(id string, messageType messages.ServerMessageType, payload interface{}) {
	recipents := make([]string, 0)
	for _, p := range r.players {
		if p.profile.Id != id {
			recipents = append(recipents, p.profile.Id)
		}
	}

	r.snd <- messages.ServerMessage{
		To:      recipents,
		Type:    messageType,
		Payload: payload,
	}
}

func (r *Room) cleanup() {
	close(r.snd)
	fblog.Info(component.Room, "room stopping", "id", r.id)
}
