package networking

import (
	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/feelbeatapp/feelbeatserver/internal/fblog"
)

type HubClient interface {
	Send([]byte)
	// Closes with notifing client
	Close()
	// Closes immediately without sending any closing message
	CloseNow()
}

type Hub struct {
	clients    map[HubClient]bool
	broadcast  chan ClientMessage
	register   chan HubClient
	unregister chan HubClient
	exit       chan bool
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[HubClient]bool),
		broadcast:  make(chan ClientMessage),
		register:   make(chan HubClient),
		unregister: make(chan HubClient),
		exit:       make(chan bool),
	}
}

func (h *Hub) Run() {
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

func (h *Hub) RegisterClient(client HubClient) {
	h.register <- client
}

func (h *Hub) Broadcast(message ClientMessage) {
	h.broadcast <- message
}

func (h *Hub) UnregisterClient(client HubClient) {
	h.unregister <- client
}

func (h *Hub) Stop() {
	h.exit <- true
}
