package networking

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func ServeWebsockets(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Origin: %v", r.Header)
		log.Fatalf("Failed to upgrade connection: %v", err)
		return
	}

	client := newClient(conn, hub.broadcast)
	hub.register <- client

	go client.readLoop()
	go client.sendLoop()
}
