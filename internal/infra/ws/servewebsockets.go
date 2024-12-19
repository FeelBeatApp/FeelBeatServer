package ws

import (
	"net/http"
	"os"

	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func ServeWebsockets(hub *BasicHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	fblog.Info(component.WebSocket, "received new connection", "ip", r.RemoteAddr)

	if err != nil {
		fblog.Error(component.WebSocket, "failed to upgrade connection", "ip", r.RemoteAddr)
		os.Exit(1)
		return
	}

	client := newClient(conn, hub)
	hub.RegisterClient(client)

	go client.readLoop()
	go client.sendLoop()
}
