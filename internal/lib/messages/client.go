package messages

import "github.com/feelbeatapp/feelbeatserver/internal/lib"

type ClientMessageType string

const (
	JoiningPlayer  = "JOIN"
	LeavingPlayer  = "LEAVE"
	SettingsUpdate = "SETTINGS_UPDATE"
	ReadyStatus    = "READY_STATUS"
	GuessSong      = "GUESS_SONG"
)

type ClientMessage struct {
	Type    ClientMessageType `json:"type"`
	From    string            `json:"-"`
	Payload interface{}       `json:"payload"`
}

type JoiningPlayerPayload struct {
	User lib.UserProfile
}

type SettingsUpdatePayload struct {
	Token    string           `json:"token"`
	Settings lib.RoomSettings `json:"settings"`
}

type GuessSongPayload struct {
	Id     string `json:"id"`
	Points int    `json:"points"`
}
