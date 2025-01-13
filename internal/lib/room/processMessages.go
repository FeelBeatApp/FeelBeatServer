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
		case messages.SettingsUpdate:
			payload, ok := message.Payload.(messages.SettingsUpdatePayload)
			if !ok {
				logIncorrectPayload("Incorrect payload in settings update", message.Payload, message.From)
			} else {
				r.updateSettings(message.From, payload)
			}
		case messages.ReadyStatus:
			if ready, ok := message.Payload.(bool); ok {
				r.updateReady(message.From, ready)
			} else {
				logIncorrectPayload("Incorrect paylod in ready status", message.Payload, message.From)
			}
		case messages.GuessSong:
			if payload, ok := message.Payload.(messages.GuessSongPayload); ok {
				r.verifyGuess(message.From, payload.Id, payload.Points)
			} else {
				logIncorrectPayload("Incorrect paylod in ready status", message.Payload, message.From)
			}
		default:
			fblog.Warn(component.Room, "Received unexpected message", "room", r.id, "from", message.From, "type", message.Type, "payload", message.Payload)
		}
	}
}

func logIncorrectPayload(text string, payload interface{}, from string) {
	fblog.Error(component.Room, text, "payload", payload, "from", from)
}
