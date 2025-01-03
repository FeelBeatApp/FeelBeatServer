package room

type Room struct {
	id       string
	ownerId  string
	settings RoomSettings
}

func NewRoom(id string, ownerId string, settings RoomSettings) Room {
	return Room{
		id:       id,
		settings: settings,
		ownerId:  ownerId,
	}
}

func (r Room) Id() string {
	return r.id
}
