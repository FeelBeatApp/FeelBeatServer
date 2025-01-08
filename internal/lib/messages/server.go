package messages

type ServerMessageType string

type ServerMessage struct {
	To      []string
	Payload interface{}
}

const KickUser = "KICK"
