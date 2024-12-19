package ws

type ClientMessage struct {
	From    HubClient
	Payload []byte
}
