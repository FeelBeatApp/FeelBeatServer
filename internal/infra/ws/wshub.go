package ws

import (
	"encoding/json"
	"fmt"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/feelbeaterror"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/messages"
)

type WSHub struct {
	clients    map[string]WSHubUser
	snd        <-chan messages.ServerMessage
	rcv        chan messages.ClientMessage
	register   chan messages.UserClient
	unregister chan string
	closed     chan bool
}

type WSHubUser struct {
	userClient messages.UserClient
	snd        chan []byte
}

func NewWSHub() messages.Hub {
	return &WSHub{
		clients:    make(map[string]WSHubUser),
		rcv:        make(chan messages.ClientMessage),
		register:   make(chan messages.UserClient),
		unregister: make(chan string),
		closed:     make(chan bool),
	}
}

func (h *WSHub) Run(snd <-chan messages.ServerMessage) <-chan messages.ClientMessage {
	h.snd = snd
	go h.run()
	return h.rcv
}

func (h *WSHub) Register(user messages.UserClient) error {
	select {
	case h.register <- user:
		return nil
	case <-h.closed:
		return fmt.Errorf("Hub is already closed")
	}
}

func (h *WSHub) run() {
	defer func() {
		close(h.register)
		close(h.unregister)

		for msg := range h.snd {
			h.sendMessage(msg)
		}

		close(h.rcv)

		for _, client := range h.clients {
			close(client.snd)
			delete(h.clients, client.userClient.User.Profile.Id)
		}
	}()

	for {
		select {
		case usersocket := <-h.register:
			hubUser := WSHubUser{
				userClient: usersocket,
				snd:        make(chan []byte),
			}
			h.clients[usersocket.User.Profile.Id] = hubUser
			go usersocket.Client.Run(hubUser.snd)
			go h.passMessages(hubUser.userClient.User.Profile.Id, hubUser.userClient.Client.ReceiveChannel())

			h.rcv <- messages.ClientMessage{
				Type: messages.JoiningPlayer,
				From: usersocket.User.Profile.Id,
				Payload: messages.JoiningPlayerPayload{
					User: usersocket.User.Profile,
				},
			}
		case userId := <-h.unregister:
			h.rcv <- messages.ClientMessage{
				Type: messages.LeavingPlayer,
				From: userId,
			}
			close(h.clients[userId].snd)
			delete(h.clients, userId)

		case message, ok := <-h.snd:
			if !ok {
				close(h.closed)
				return
			}
			h.sendMessage(message)
		}
	}
}

func (h *WSHub) passMessages(from string, rcv <-chan []byte) {
	for bytes := range rcv {
		var message messages.ClientMessage
		err := json.Unmarshal(bytes, &message)
		if err != nil {
			fblog.Error(component.Hub, "Failed to parse client message", "err", err)
		}
		message.From = from
		fmt.Println("Check")
		fmt.Println(message)
		fmt.Println(from)

		h.rcv <- message
	}

	h.unregister <- from
}

func (h *WSHub) sendMessage(message messages.ServerMessage) {
	if msg, ok := message.Payload.(string); ok {
		if msg == messages.KickUser {
			for _, id := range message.To {
				close(h.clients[id].snd)
			}
			return
		}
	}

	bytes, err := json.Marshal(message.Payload)
	if err != nil {
		bytes = []byte(feelbeaterror.EncodingMessageFailed)
	}

	for _, id := range message.To {
		h.clients[id].snd <- bytes
	}
}
