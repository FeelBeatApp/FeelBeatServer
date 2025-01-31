package messages

import (
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
)

type ServerMessageType string

type ServerMessage struct {
	To      []string
	Type    ServerMessageType
	Payload interface{}
}

func (m ServerMessage) ToUnit() ServerMessageUnit {
	return ServerMessageUnit{
		Type:    m.Type,
		Payload: m.Payload,
	}
}

type ServerMessageUnit struct {
	Type    ServerMessageType `json:"type"`
	Payload interface{}       `json:"payload"`
}

const (
	InitialMessage = "INITIAL"
	NewPlayer      = "NEW_PLAYER"
	PlayerLeft     = "PLAYER_LEFT"
	ServerError    = "SERVER_ERROR"
	RoomStage      = "ROOM_STAGE"
	PlayerReady    = "PLAYER_READY"
	PlaySong       = "PLAY_SONG"
	PlayerGuess    = "PLAYER_GUESS"
	CorrectSong    = "CORRECT_SONG"
	EndGame        = "END_GAME"
)

type InitialGameState struct {
	Id       string            `json:"id"`
	Me       string            `json:"me"`
	Admin    string            `json:"admin"`
	Playlist PlaylistState     `json:"playlist"`
	Players  []lib.UserProfile `json:"players"`
	Settings lib.RoomSettings  `json:"settings"`
	ReadyMap map[string]bool   `json:"readyMap"`
}

type PlaylistState struct {
	Name     string      `json:"name"`
	ImageUrl string      `json:"imageUrl"`
	Songs    []SongState `json:"songs"`
}

type SongState struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	ImageUrl    string `json:"imageUrl"`
	DurationSec int    `json:"durationSec"`
}

type PlayerLeftPayload struct {
	Left  string `json:"left"`
	Admin string `json:"admin"`
}

type PlayerReadyPayload struct {
	Player string `json:"player"`
	Ready  bool   `json:"ready"`
}

type PlaySongPayload struct {
	Url       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
	Duration  int64  `json:"duration"`
}

type PlayerGuessPayload struct {
	Correct  bool   `json:"correct"`
	Points   int    `json:"points"`
	PlayerId string `json:"playerId"`
	SongId   string `json:"songId"`
}

type EndGamePayload struct {
	Results []PlayerResult `json:"results"`
}

type PlayerResult struct {
	Profile lib.UserProfile `json:"profile"`
	Points  int             `json:"points"`
}
