package ws

import (
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
)

type BasicHub struct {
	clients    map[HubClient]bool
	broadcast  chan ClientMessage
	register   chan HubClient
	unregister chan HubClient
	exit       chan bool
}

func NewHub() *BasicHub {
	return &BasicHub{
		clients:    make(map[HubClient]bool),
		broadcast:  make(chan ClientMessage),
		register:   make(chan HubClient),
		unregister: make(chan HubClient),
		exit:       make(chan bool),
	}
}

func (h *BasicHub) Run() {
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

func (h *BasicHub) RegisterClient(client HubClient) {
	h.register <- client
}

func (h *BasicHub) Broadcast(message ClientMessage) {
	h.broadcast <- message
}

func (h *BasicHub) UnregisterClient(client HubClient) {
	h.unregister <- client
}

func (h *BasicHub) Stop() {
	h.exit <- true
}
