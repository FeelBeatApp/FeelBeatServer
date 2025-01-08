package roomapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
)

type fetchGamesResponse struct {
	Rooms []responseRoom `json:"rooms"`
}

type responseRoom struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Players    int    `json:"players"`
	MaxPlayers int    `json:"maxPlayers"`
	ImageUrl   string `json:"imageUrl"`
}

func (r RoomApi) fetchRoomsHandler(user auth.User, res http.ResponseWriter, req *http.Request) {
	rooms := r.roomRepo.GetAllRooms()
	formatted := fetchGamesResponse{
		Rooms: make([]responseRoom, 0),
	}

	for _, room := range rooms {
		formatted.Rooms = append(formatted.Rooms, responseRoom{
			Id:         room.Id(),
			Name:       room.Name(),
			Players:    len(room.PlayerProfiles()),
			MaxPlayers: room.Settings().MaxPlayers,
			ImageUrl:   room.ImageUrl(),
		})
	}

	resJson, err := json.Marshal(formatted)
	if err != nil {
		api.LogApiError("Couldn't encode response", err, user.Profile.Id, req)
		return
	}
	err = api.SendJsonResponse(&res, resJson)
	if err != nil {
		api.LogApiError("Couldn't write response", err, user.Profile.Id, req)
		return
	}

	api.LogApiCall(user.Profile.Id, req)
}

func (r RoomApi) ServeFetchRooms(baseUrl string, authWrapper auth.AuthWrapper) {
	http.HandleFunc(fmt.Sprintf("%s/rooms", baseUrl), authWrapper(r.fetchRoomsHandler))
}
