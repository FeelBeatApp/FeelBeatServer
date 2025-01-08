package messages

type Hub interface {
	Run(<-chan ServerMessage) <-chan ClientMessage
	Register(UserClient) error
}

type HubClient interface {
	Run(<-chan []byte)
	ReceiveChannel() <-chan []byte
}
