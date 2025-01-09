package room

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/messages"
)

func (r *Room) processMessages() {
	for message := range r.rcv {
		switch message.Type {
		case messages.JoiningPlayer:
			payload, ok := message.Payload.(messages.JoiningPlayerPayload)
			if !ok {
				logIncorrectPayload("Incorrect payload in player joining", message.Payload, message.From)
			} else {
				r.addPlayer(payload.User)
			}
		case messages.LeavingPlayer:
			r.removePlayer(message.From)
		default:
			fblog.Warn(component.Room, "Received unexpected message", "room", r.id, "from", message.From, "type", message.Type, "payload", message.Payload)
		}
	}
}

func logIncorrectPayload(text string, payload interface{}, from string) {
	fblog.Error(component.Room, text, "payload", payload, "from", from)
}
