package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/buger/jsonparser"
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
	isOpen     bool
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
		isOpen:     true,
	}
}

func (h *WSHub) Run(snd <-chan messages.ServerMessage) <-chan messages.ClientMessage {
	h.snd = snd
	go h.run()
	return h.rcv
}

func (h *WSHub) Register(user messages.UserClient) error {
	if !h.isOpen {
		return fmt.Errorf("Hub is not open")
	}

	h.register <- user
	return nil
}

func (h *WSHub) run() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	defer func() {
		h.isOpen = false
		cancel()
		wg.Wait()

		for _, c := range h.clients {
			close(c.snd)
		}

		close(h.register)
		close(h.unregister)
		close(h.rcv)

		fblog.Info(component.Hub, "hub closed")
	}()

	for {
		select {
		case usersocket := <-h.register:
			hubUser := WSHubUser{
				userClient: usersocket,
				snd:        make(chan []byte),
			}
			h.clients[usersocket.User.Profile.Id] = hubUser
			wg.Add(1)
			go usersocket.Client.Run(hubUser.snd)
			go h.passMessages(ctx, &wg, hubUser.userClient.User.Profile.Id, hubUser.userClient.Client.ReceiveChannel())

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

			if c, ok := h.clients[userId]; ok {
				close(c.snd)
			}
			delete(h.clients, userId)

		case message, ok := <-h.snd:
			if !ok {
				return
			}
			h.sendMessage(message)
		}
	}
}

func (h *WSHub) passMessages(ctx context.Context, wg *sync.WaitGroup, from string, rcv <-chan []byte) {
	defer func() {
		wg.Done()
	}()

	for {
		select {
		case bytes, ok := <-rcv:
			if !ok {
				h.unregister <- from
				return
			}

			msgType, err := jsonparser.GetString(bytes, "type")
			if err != nil {
				fblog.Error(component.Hub, "Failed to parse client message", "err", err)
				continue
			}
			payload, _, _, err := jsonparser.Get(bytes, "payload")
			if err != nil {
				fblog.Error(component.Hub, "Failed to parse client message", "err", err)
				continue
			}

			h.rcv <- decodeMessage(from, msgType, payload)
		case <-ctx.Done():
			return
		}
	}
}

func decodeMessage(from string, msgType string, payload []byte) messages.ClientMessage {
	var settingsUpdate messages.SettingsUpdatePayload

	var err error
	var result interface{}
	switch msgType {
	case messages.SettingsUpdate:
		err = json.Unmarshal(payload, &settingsUpdate)
		result = settingsUpdate
	}

	if err != nil {
		return messages.ClientMessage{
			From: from,
			Type: messages.ClientMessageType(msgType),
		}
	} else {
		return messages.ClientMessage{
			From:    from,
			Type:    messages.ClientMessageType(msgType),
			Payload: result,
		}
	}
}

func (h *WSHub) sendMessage(message messages.ServerMessage) {
	bytes, err := json.Marshal(message.ToUnit())
	if err != nil {
		bytes = []byte(feelbeaterror.EncodingMessageFailed)
	}

	fblog.Info(component.Hub, "sending message", "type", message.Type, "to", message.To)

	for _, id := range message.To {
		h.clients[id].snd <- bytes
	}
}
