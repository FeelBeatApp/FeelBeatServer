package messages

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
)

type UserClient struct {
	Client HubClient
	User   auth.User
}

func NewUserSocket(client HubClient, user auth.User) UserClient {
	return UserClient{
		Client: client,
		User:   user,
	}
}
