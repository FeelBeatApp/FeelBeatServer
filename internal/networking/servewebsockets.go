package networking

import (
	"net/http"
	"os"

	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/feelbeatapp/feelbeatserver/internal/fblog"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func ServeWebsockets(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	fblog.Info(component.WebSocket, "received new connection", "ip", r.RemoteAddr)

	if err != nil {
		fblog.Error(component.WebSocket, "failed to upgrade connection", "ip", r.RemoteAddr)
		os.Exit(1)
		return
	}

	client := newClient(conn, hub.broadcast, hub.unregister)
	hub.register <- client

	go client.readLoop()
	go client.sendLoop()
}
