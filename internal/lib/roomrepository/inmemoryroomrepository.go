package roomrepository

import (
	"sync"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/messages"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/room"
	"github.com/google/uuid"
)

type InMemoryRoomRepository struct {
	createHub func() messages.Hub
	spotify   lib.SpotifyApi
	rooms     map[string]*room.Room
	m         sync.RWMutex
}

func NewInMemoryRoomRepository(spotify lib.SpotifyApi, createHub func() messages.Hub) *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		createHub: createHub,
		spotify:   spotify,
		rooms:     make(map[string]*room.Room),
	}
}

func (r *InMemoryRoomRepository) CreateRoom(user auth.User, settings lib.RoomSettings) (string, error) {
	playlistData, err := r.spotify.FetchPlaylistData(settings.PlaylistId, user.Token)
	if err != nil {
		return "", err
	}

	newRoom := room.NewRoom(uuid.NewString(), playlistData, user.Profile, settings, r.createHub(), r.spotify, func(room *room.Room) {
		r.m.Lock()
		delete(r.rooms, room.Id())
		r.m.Unlock()
		fblog.Info(component.RoomRepository, "removed room", "id", room.Id())
	})
	r.m.Lock()
	r.rooms[newRoom.Id()] = newRoom
	r.m.Unlock()

	newRoom.Start()

	fblog.Info(component.RoomRepository, "room created and started", "id", newRoom.Id(), "room count", len(r.rooms))

	return newRoom.Id(), nil
}

func (r *InMemoryRoomRepository) GetAllRooms() []*room.Room {
	result := make([]*room.Room, 0)
	r.m.RLock()
	for _, room := range r.rooms {
		result = append(result, room)
	}
	r.m.RUnlock()

	return result
}

func (r *InMemoryRoomRepository) GetRoom(id string) *room.Room {
	defer r.m.RUnlock()
	r.m.RLock()
	return r.rooms[id]
}
