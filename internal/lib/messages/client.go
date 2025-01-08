package messages

import "github.com/feelbeatapp/feelbeatserver/internal/lib"

type ClientMessageType string

const (
	JoiningPlayer = "JOIN"
	LeavingPlayer = "LEAVE"
)

type ClientMessage struct {
	Type    ClientMessageType `json:"type"`
	From    string
	Payload interface{} `json:"payload"`
}

type JoiningPlayerPayload struct {
	User lib.UserProfile
}
