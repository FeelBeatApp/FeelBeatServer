package lib

import "time"

type SongDetails struct {
	Title    string
	Artist   string
	Duration time.Duration
}

type Song struct {
	Id      string
	Details SongDetails
}

type PlaylistData struct {
	Name     string
	ImageUrl string
	Songs    []Song
}

type UserProfile struct {
	Id       string
	Name     string
	ImageUrl string
}
