package roomapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/feelbeaterror"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
)

type createGameResponse struct {
	RoomId string `json:"roomId"`
}

func (r RoomApi) createGameHandler(userId string, res http.ResponseWriter, req *http.Request) {
	var payload room.RoomSettings
	err := api.ParseBody(req.Body, &payload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fblog.Error(component.Api, "Error: ", "err", err)
		return
	}

	roomId, err := r.roomRepo.CreateRoom(userId, payload)
	if err != nil {
		var fbError *feelbeaterror.FeelBeatError
		if errors.As(err, &fbError) {
			http.Error(res, string(fbError.UserMessage), feelbeaterror.StatusCode(fbError.UserMessage))
		} else {
			http.Error(res, feelbeaterror.Default, feelbeaterror.StatusCode(feelbeaterror.Default))
		}

		api.LogApiError("Create room failed", err, userId, req)
		return
	}

	resJson, err := json.Marshal(createGameResponse{
		RoomId: roomId,
	})
	if err != nil {
		api.LogApiError("Couldn't encode response", err, userId, req)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resJson)
	if err != nil {
		api.LogApiError("Couldn't write response", err, userId, req)
		return
	}

	api.LogApiCall(userId, req)
}

func (r RoomApi) ServeCreateGame(baseUrl string, authWrapper auth.AuthWrapper) {
	http.HandleFunc(fmt.Sprintf("%s/create", baseUrl), authWrapper(r.createGameHandler))
}