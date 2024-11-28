package networking

import (
	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/feelbeatapp/feelbeatserver/internal/fblog"
	"github.com/gorilla/websocket"
)

// TODO: Implement ping pong

const DEFAULT_OUT_BUFFER_SIZE = 256

type Client struct {
	broadcast  chan<- ClientMessage
	unregister chan<- *Client
	conn       *websocket.Conn
	send       chan []byte
}

func newClient(conn *websocket.Conn, hubChannel chan ClientMessage, unregister chan *Client) *Client {
	return &Client{
		broadcast:  hubChannel,
		unregister: unregister,
		conn:       conn,
		send:       make(chan []byte, DEFAULT_OUT_BUFFER_SIZE),
	}
}

func (c *Client) readLoop() {
	defer c.conn.Close()

	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fblog.Info(component.Client, "client closed connection", "ip", c.conn.RemoteAddr())
			} else {
				fblog.Error(component.Client, "Received unexpected error from client", "err", err)
			}
			c.unregister <- c
			break
		}

		switch msgType {
		case websocket.TextMessage:
			c.broadcast <- ClientMessage{
				from:    c,
				payload: message,
			}
		default:
			fblog.Warn(component.Client, "ignoring message", "type", msgType, "msg", message)
		}
	}
}

func (c *Client) sendLoop() {
	defer c.conn.Close()

	for {
		message, ok := <-c.send
		if !ok {
			fblog.Info(component.Client, "server is closing connection", "with", c.conn.RemoteAddr())
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			fblog.Error(component.Client, "error when sending message", "msg", message)
		}
	}
}

func (c *Client) Close() {
	err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		fblog.Error(component.Client, "couldn't properly close connection", "err", err)
	}

	close(c.send)
}
