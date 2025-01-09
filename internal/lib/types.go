package lib

import "time"

type SongDetails struct {
	Title    string        `json:"title"`
	Artist   string        `json:"artist"`
	ImageUrl string        `json:"imageUrl"`
	Duration time.Duration `json:"duration"`
}

type Song struct {
	Id      string      `json:"id"`
	Details SongDetails `json:"details"`
}

type PlaylistData struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	Songs    []Song `json:"songs"`
}

type UserProfile struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

type RoomSettings struct {
	MaxPlayers            int    `json:"maxPlayers"`
	TurnCount             int    `json:"turnCount"`
	TimePenaltyPerSecond  int    `json:"timePenaltyPerSecond"`
	BasePoints            int    `json:"basePoints"`
	IncorrectGuessPenalty int    `json:"incorrectGuessPenalty"`
	PlaylistId            string `json:"playlistId"`
}
