package networking_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/networking"
	"github.com/stretchr/testify/assert"
)

type FakeClient struct {
	payloads [][]byte
	closed   bool
}

func newFakeClient() *FakeClient {
	return &FakeClient{
		payloads: make([][]byte, 0),
		closed:   false,
	}
}

func (c *FakeClient) Send(payload []byte) {
	fmt.Println("Let's go I'm really sending sht")
	c.payloads = append(c.payloads, payload)
	fmt.Println(c.payloads)
}

func (c *FakeClient) CloseNow() {
	c.closed = true
}

func (c *FakeClient) Close() {
	c.closed = true
}

const testMessage = "hi there"

func TestHubBroadcastsMessages(t *testing.T) {
	assert := assert.New(t)
	hub := networking.NewHub()

	go hub.Run()

	clients := make([]*FakeClient, 5)
	for i := range clients {
		clients[i] = newFakeClient()
		hub.RegisterClient(clients[i])
	}

	hub.Broadcast(networking.ClientMessage{
		From:    clients[0],
		Payload: []byte(testMessage),
	})

	assert.Equal(0, len(clients[0].payloads))

	time.Sleep(time.Millisecond * 1)

	for i := 1; i < len(clients); i++ {
		assert.NotEmpty(clients[i].payloads)
		assert.Contains(clients[i].payloads, []byte(testMessage))
	}

	hub.Stop()
}
