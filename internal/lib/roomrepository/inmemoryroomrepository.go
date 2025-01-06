package roomrepository

import (
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/validation"
	"github.com/google/uuid"
)

type SpotifyApi interface {
	FetchPlaylistSongs(playlistId string, token string) ([]lib.Song, error)
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

// TODO: Implement fetching playlist details from spotify
func (r InMemoryRoomRepository) CreateRoom(ownderId string, settings room.RoomSettings, token string) (string, error) {
	err := validation.ValidateRoomSettings(settings)
	if err != nil {
		return "", err
	}

	songs, err := r.spotify.FetchPlaylistSongs(settings.PlaylistId, token)
	if err != nil {
		return "", err
	}

	newRoom := room.NewRoom(uuid.NewString(), ownderId, settings, songs)
	r.rooms[newRoom.Id()] = newRoom

	return newRoom.Id(), nil
}
