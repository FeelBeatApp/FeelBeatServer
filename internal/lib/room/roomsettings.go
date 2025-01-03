package room

type RoomSettings struct {
	MaxPlayers            int    `json:"maxPlayers"`
	TurnCount             int    `json:"turnCount"`
	TimePenaltyPerSecond  int    `json:"timePenaltyPerSecond"`
	BasePoints            int    `json:"basePoints"`
	IncorrectGuessPenalty int    `json:"incorrectGuessPenalty"`
	PlaylistId            string `json:"playlistId"`
}
