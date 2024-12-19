package ws

type Hub interface {
	RegisterClient(HubClient)
	UnregisterClient(HubClient)
	Broadcast(ClientMessage)
}

type HubClient interface {
	Send([]byte)
	// Closes with notifing client
	Close()
	// Closes immediately without sending any closing message
	CloseNow()
}
