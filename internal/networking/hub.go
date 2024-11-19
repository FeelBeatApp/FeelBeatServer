package networking

import "log/slog"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			slog.Info("Client registered")
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
			slog.Info("Client left")
		case message := <-h.broadcast:
			slog.Info("Broadcasting", "msg", string(message))
			for client := range h.clients {
				select {

				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}
