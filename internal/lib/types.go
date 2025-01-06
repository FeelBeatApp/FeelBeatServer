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

type Playlist struct {
	Id    string
	Songs []Song
}
