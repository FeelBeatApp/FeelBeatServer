package networking

type ClientMessage struct {
	From    HubClient
	Payload []byte
}
