package networking

import (
	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/feelbeatapp/feelbeatserver/internal/fblog"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan ClientMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan ClientMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fblog.Info(component.Hub, "new client registered")
		case client := <-h.unregister:
			delete(h.clients, client)
			fblog.Info(component.Client, "client unregistered")
		case message := <-h.broadcast:
			for client := range h.clients {
				if client != message.from {
					client.send <- message.payload
				}
			}
		}
	}
}
