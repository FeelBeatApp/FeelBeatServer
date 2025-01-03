package ws

import (
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
)

type WSHub struct {
	clients    map[HubClient]bool
	broadcast  chan ClientMessage
	register   chan HubClient
	unregister chan HubClient
	exit       chan bool
}

func NewHub() *WSHub {
	return &WSHub{
		clients:    make(map[HubClient]bool),
		broadcast:  make(chan ClientMessage),
		register:   make(chan HubClient),
		unregister: make(chan HubClient),
		exit:       make(chan bool),
	}
}

func (h *WSHub) Run() {
	defer func() {
		for c := range h.clients {
			c.Close()
		}
		close(h.broadcast)
		close(h.register)
		close(h.unregister)
		close(h.exit)
	}()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fblog.Info(component.Hub, "new client registered")
		case client := <-h.unregister:
			delete(h.clients, client)
			client.CloseNow()
			fblog.Info(component.Client, "client unregistered")
		case message := <-h.broadcast:
			for client := range h.clients {
				if client != message.From {
					client.Send(message.Payload)
				}
			}
		case <-h.exit:
			break
		}
	}
}

func (h *WSHub) RegisterClient(client HubClient) {
	h.register <- client
}

func (h *WSHub) Broadcast(message ClientMessage) {
	h.broadcast <- message
}

func (h *WSHub) UnregisterClient(client HubClient) {
	h.unregister <- client
}

func (h *WSHub) Stop() {
	h.exit <- true
}
