package networking

type ClientMessage struct {
	from    *Client
	payload []byte
}
