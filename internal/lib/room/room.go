package room

import (
	"errors"
	"math/rand/v2"
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/audioprovider"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/feelbeaterror"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/messages"
)

const songPlayDelay = time.Second * 3

type Room struct {
	id         string
	playlist   lib.PlaylistData
	owner      lib.UserProfile
	settings   lib.RoomSettings
	players    map[string]*Player
	readyMap   map[string]bool
	stage      RoomStage
	hub        messages.Hub
	snd        chan messages.ServerMessage
	rcv        <-chan messages.ClientMessage
	onCleanup  func(*Room)
	spotifyApi lib.SpotifyApi
	audio      audioprovider.AudioProvider
	game       gameState
}

type Player struct {
	profile lib.UserProfile
	points  int
}

type gameState struct {
	turn       int
	pickedSong *lib.Song
}

func NewRoom(id string,
	playlist lib.PlaylistData,
	owner lib.UserProfile,
	settings lib.RoomSettings,
	hub messages.Hub,
	spotifyApi lib.SpotifyApi,
	audio audioprovider.AudioProvider,
	onCleanup func(*Room),
) *Room {
	return &Room{
		id:         id,
		playlist:   playlist,
		owner:      owner,
		settings:   settings,
		players:    make(map[string]*Player),
		readyMap:   make(map[string]bool),
		stage:      LobbyStage,
		hub:        hub,
		snd:        make(chan messages.ServerMessage),
		spotifyApi: spotifyApi,
		audio:      audio,
		onCleanup:  onCleanup,
		game: gameState{
			turn:       0,
			pickedSong: nil,
		},
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

func (r *Room) Stage() RoomStage {
	return r.stage
}

func (r *Room) Start() {
	r.rcv = r.hub.Run(r.snd)
	go r.processMessages()
}

func (r *Room) Hub() messages.Hub {
	return r.hub
}

func (r *Room) addPlayer(profile lib.UserProfile) {
	if r.stage != LobbyStage {
		return
	}

	r.players[profile.Id] = &Player{
		profile: profile,
		points:  0,
	}

	fblog.Info(component.Room, "new player", "roomId", r.id, "userId", profile.Id)

	r.snd <- r.packIntialState(profile.Id)
	r.sendToAllExcept(profile.Id, messages.NewPlayer, profile)
}

func (r *Room) packIntialState(me string) messages.ServerMessage {
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

	return messages.ServerMessage{
		To:   []string{me},
		Type: messages.InitialMessage,
		Payload: messages.InitialGameState{
			Id:    r.id,
			Me:    me,
			Admin: r.owner.Id,
			Playlist: messages.PlaylistState{
				Name:     r.Name(),
				ImageUrl: r.ImageUrl(),
				Songs:    packedSongs,
			},
			Players:  playerProfiles,
			ReadyMap: r.readyMap,
			Settings: r.settings,
		},
	}

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

	if r.stage == LobbyStage && r.allReady() {
		r.broadcastRoomStage(GameStage)
		r.provideSong()
	}
}

func (r *Room) updateSettings(from string, settingsPayload messages.SettingsUpdatePayload) {
	if settingsPayload.Settings.MaxPlayers < len(r.players) {
		return
	}

	ok := true
	if settingsPayload.Settings.PlaylistId != r.settings.PlaylistId {
		playlistData, err := r.spotifyApi.FetchPlaylistData(settingsPayload.Settings.PlaylistId, settingsPayload.Token)
		if err != nil {
			ok = false
			fblog.Error(component.Room, "failed to change playlist", "roomId", r.id, "err", err)
			var fbErr *feelbeaterror.FeelBeatError
			if errors.As(err, &fbErr) {

				r.snd <- messages.ServerMessage{
					Type:    messages.ServerError,
					To:      []string{from},
					Payload: fbErr.UserMessage,
				}
			}
		}

		r.playlist = playlistData

	}

	if ok {
		r.settings = settingsPayload.Settings
		fblog.Info(component.Room, "settings updated", "roomId", r.id, "settings", r.settings)

		for _, player := range r.players {
			r.snd <- r.packIntialState(player.profile.Id)
		}
	} else {
		r.snd <- r.packIntialState(from)
	}

}

func (r *Room) updateReady(from string, ready bool) {
	r.readyMap[from] = ready

	if r.allReady() {
		r.broadcastRoomStage(GameStage)
		r.provideSong()
	} else {
		r.sendToAllExcept(from, messages.PlayerReady, messages.PlayerReadyPayload{
			Player: from,
			Ready:  ready,
		})
	}
}

func (r *Room) allReady() bool {
	allReady := true
	for _, p := range r.players {
		allReady = allReady && r.readyMap[p.profile.Id]
	}
	return allReady
}

func (r *Room) broadcastRoomStage(stage RoomStage) {
	r.stage = stage
	recipents := make([]string, 0)
	for _, p := range r.players {
		recipents = append(recipents, p.profile.Id)
	}

	r.snd <- messages.ServerMessage{
		To:      recipents,
		Type:    messages.RoomStage,
		Payload: stage,
	}
}

func (r *Room) provideSong() {
	var url string
	for {
		r.game.pickedSong = &r.playlist.Songs[rand.IntN(len(r.playlist.Songs))]
		u, err := r.audio.GetUrl(r.game.pickedSong.Details)
		url = u
		if err == nil {
			break
		}
	}

	recipents := make([]string, 0)
	for _, p := range r.players {
		recipents = append(recipents, p.profile.Id)
	}

	start := time.Now().Add(songPlayDelay)

	r.snd <- messages.ServerMessage{
		To:   recipents,
		Type: messages.PlaySong,
		Payload: messages.PlaySongPayload{
			Url:       url,
			Timestamp: start.Unix(),
			Duration:  r.game.pickedSong.Details.Duration.Milliseconds(),
		},
	}
}

func (r *Room) verifyGuess(from string, guess string, points int) {
	if r.stage != GameStage || r.game.pickedSong == nil {
		return
	}

	correct := false
	receivedPoints := points

	if r.game.pickedSong.Id == guess {
		r.players[from].points += points
		receivedPoints = points
		correct = true
	} else {
		r.players[from].points -= r.settings.IncorrectGuessPenalty
		receivedPoints = -r.settings.IncorrectGuessPenalty
		correct = false
	}

	r.sendToAllExcept(from, messages.PlayerGuess, messages.PlayerGuessPayload{
		Correct:  correct,
		Points:   receivedPoints,
		PlayerId: from,
		SongId:   "",
	})

	r.snd <- messages.ServerMessage{
		Type: messages.PlayerGuess,
		To:   []string{from},
		Payload: messages.PlayerGuessPayload{
			Correct:  correct,
			Points:   receivedPoints,
			PlayerId: from,
			SongId:   guess,
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
