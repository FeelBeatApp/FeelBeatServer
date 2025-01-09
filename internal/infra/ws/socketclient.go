package ws

import (
	"context"
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong mesage from the peer
	pongWait = 60 * time.Second

	// Send ping frame with this period
	pingPeriod = (pongWait * 9) / 10
)

type SocketClient struct {
	conn   *websocket.Conn
	rcv    chan []byte
	snd    <-chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}

func newSocketClient(conn *websocket.Conn) *SocketClient {
	return &SocketClient{
		conn: conn,
		rcv:  make(chan []byte),
	}
}

func (s *SocketClient) Run(snd <-chan []byte) {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.snd = snd

	go s.receivingLoop()
	go s.sendingLoop()
}

func (s *SocketClient) ReceiveChannel() <-chan []byte {
	return s.rcv
}

func (s *SocketClient) setPongDeadline() error {
	err := s.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		fblog.Error(component.Socket, "Pong dead line error: ", "err", err, "conn", s.conn)
		return err
	}

	return nil
}

func (s *SocketClient) receivingLoop() {
	defer func() {
		close(s.rcv)
		s.cancel()
	}()

	err := s.setPongDeadline()
	if err != nil {
		return
	}
	s.conn.SetPongHandler(func(string) error {
		return s.setPongDeadline()
	})

	for {
		msgType, message, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fblog.Info(component.Socket, "client closed connection", "ip", s.conn.RemoteAddr())
			} else {
				fblog.Error(component.Socket, "client connection closing due to unexpected error", "err", err, "ip", s.conn.RemoteAddr())
			}
			return
		}

		switch msgType {
		case websocket.TextMessage:
			select {
			case <-s.ctx.Done():
				return
			case s.rcv <- message:
			}
		default:
			fblog.Warn(component.Socket, "ignoring unexpected message", "type", msgType, "msg", string(message))
		}
	}
}

func (s *SocketClient) sendingLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

		for msg := range s.snd {
			err := s.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fblog.Error(component.Socket, "error during draning send channel after server closing", "with", s.conn.RemoteAddr())
			}
		}

		s.conn.Close()
		s.cancel()
	}()

	for {
		select {
		case message, ok := <-s.snd:
			if !ok {
				fblog.Info(component.Socket, "server is closing connection", "with", s.conn.RemoteAddr())
				return
			}

			err := s.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fblog.Error(component.Socket, "error when sending message", "msg", message, "err", err)
			}
		case <-ticker.C:
			err := s.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				fblog.Error(component.Socket, "couldn't ping client", "ip", s.conn.RemoteAddr())
			}
		case <-s.ctx.Done():
			return
		}
	}
}
