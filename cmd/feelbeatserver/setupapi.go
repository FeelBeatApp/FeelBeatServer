package main

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/api/roomapi"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/roomrepository"
)

const baseUrl = "/api/v1"

// TODO: Add contexts to api handler for graceful shutdown

func setupAPI(authWrapper auth.AuthWrapper, roomRepo roomrepository.RoomRepository) {
	roomApi := roomapi.NewRoomApi(roomRepo)

	handlers := []func(string, auth.AuthWrapper){roomApi.ServeCreateGame}

	fblog.Info(component.Api, "Setting up REST API", "handlers count", len(handlers))

	for _, f := range handlers {
		f(baseUrl, authWrapper)
	}
}
