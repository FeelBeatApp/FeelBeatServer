package networking

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

// TODO: Implement ping pong

const DEFAULT_OUT_BUFFER_SIZE = 256

type Client struct {
	broadcast chan<- []byte
	conn      *websocket.Conn
	send      chan []byte
}

func newClient(conn *websocket.Conn, hubChannel chan []byte) *Client {
	return &Client{
		broadcast: hubChannel,
		conn:      conn,
		send:      make(chan []byte, DEFAULT_OUT_BUFFER_SIZE),
	}
}

func (c *Client) readLoop() {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Error", "err", err)
			}
			break
		}

		c.broadcast <- message
	}
}

func (c *Client) sendLoop() {
	defer c.conn.Close()

	for {
		message, ok := <-c.send
		if !ok {
			err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				slog.Error("Error", "err", err)
			}

			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			slog.Error("Error", "err", err)
			return
		}

		_, err = w.Write(message)
		if err != nil {
			slog.Error("Error", "err", err)
		}

		if err = w.Close(); err != nil {
			slog.Error("Error", "err", err)
			return
		}
	}
}
