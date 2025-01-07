package roomrepository

import (
	"fmt"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
	"github.com/google/uuid"
)

type SpotifyApi interface {
	FetchPlaylistData(playlistId string, token string) (lib.PlaylistData, error)
}

type InMemoryRoomRepository struct {
	spotify SpotifyApi
	rooms   map[string]room.Room
}

func NewInMemoryRoomRepository(spotify SpotifyApi) InMemoryRoomRepository {
	return InMemoryRoomRepository{
		spotify: spotify,
		rooms:   make(map[string]room.Room),
	}
}

func (r InMemoryRoomRepository) CreateRoom(user auth.User, settings room.RoomSettings) (string, error) {
	playlistData, err := r.spotify.FetchPlaylistData(settings.PlaylistId, user.Token)
	if err != nil {
		return "", err
	}

	newRoom := room.NewRoom(uuid.NewString(), playlistData, user.Profile, settings)
	r.rooms[newRoom.Id()] = newRoom

	fmt.Println(newRoom)

	return newRoom.Id(), nil
}

func (r InMemoryRoomRepository) GetAllRooms() []room.Room {
	result := make([]room.Room, 0)
	for _, room := range r.rooms {
		result = append(result, room)
	}

	return result
}
