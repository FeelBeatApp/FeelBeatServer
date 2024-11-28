package networking

import (
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/feelbeatapp/feelbeatserver/internal/fblog"
	"github.com/gorilla/websocket"
)

const (
	defaultOutBufferSize = 256

	// Time allowed to read the next pong mesage from the peer
	pongWait = 60 * time.Second

	// Send ping frame with this period
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	broadcast  chan<- ClientMessage
	unregister chan<- HubClient
	conn       *websocket.Conn
	send       chan []byte
}

func newClient(conn *websocket.Conn, hubChannel chan ClientMessage, unregister chan HubClient) *Client {
	return &Client{
		broadcast:  hubChannel,
		unregister: unregister,
		conn:       conn,
		send:       make(chan []byte, defaultOutBufferSize),
	}
}

func (c *Client) setPongDeadline() {
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fblog.Error(component.Client, "Unexpected error when setting up connection", "err", err)
	}
}

func (c *Client) readLoop() {
	defer func() {
		c.unregister <- c
		c.conn.Close()
	}()

	c.setPongDeadline()
	c.conn.SetPongHandler(func(string) error {
		c.setPongDeadline()
		return nil
	})

	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fblog.Info(component.Client, "client closed connection", "ip", c.conn.RemoteAddr())
			} else {
				fblog.Error(component.Client, "Received unexpected error from client", "err", err)
			}
			break
		}

		switch msgType {
		case websocket.TextMessage:
			c.broadcast <- ClientMessage{
				From:    c,
				Payload: message,
			}
		default:
			fblog.Warn(component.Client, "ignoring message", "type", msgType, "msg", message)
		}
	}
}

func (c *Client) sendLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				fblog.Info(component.Client, "server is closing connection", "with", c.conn.RemoteAddr())
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				fblog.Error(component.Client, "error when sending message", "msg", message)
			}
		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				fblog.Error(component.Client, "couldn't ping client", "ip", c.conn.RemoteAddr())
			}
		}

	}
}

func (c *Client) Send(payload []byte) {
	c.send <- payload
}

func (c *Client) Close() {
	err := c.conn.WriteMessage(websocket.CloseMessage, nil)
	if err != nil {
		fblog.Error(component.Client, "couldn't close connection gracefully", "err", err)
		c.conn.Close()
	}
}

func (c *Client) CloseNow() {
	close(c.send)
}
