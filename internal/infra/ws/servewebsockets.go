package ws

import (
	"fmt"
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/feelbeaterror"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/messages"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func (w WSHandler) websocketHandler(user auth.User, res http.ResponseWriter, req *http.Request) {
	roomId := req.PathValue("id")
	room := w.roomRepo.GetRoom(roomId)
	if room == nil {
		http.Error(res, feelbeaterror.RoomNotFound, feelbeaterror.StatusCode(feelbeaterror.RoomNotFound))
		api.LogApiError("User tried to connect to non existing room", nil, user.Profile.Id, req)
		return
	}
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		http.Error(res, feelbeaterror.Default, feelbeaterror.StatusCode(feelbeaterror.Default))
		api.LogApiError("Upgrading websocket connection failed", err, user.Profile.Id, req)
		return
	}

	userSocket := messages.NewUserSocket(newSocketClient(conn), user)
	err = room.Hub().Register(userSocket)
	if err != nil {
		http.Error(res, feelbeaterror.RoomNotFound, feelbeaterror.StatusCode(feelbeaterror.RoomNotFound))
		api.LogApiError("Registering socket failed, hub is closed", err, user.Profile.Id, req)
		return
	}

	api.LogApiCall(user.Profile.Id, req)
}

func (w WSHandler) ServeWebsockets(basePath string, authWrapper auth.AuthWrapper) {
	http.HandleFunc(fmt.Sprintf("%s/{id}", basePath), authWrapper(w.websocketHandler))
}
