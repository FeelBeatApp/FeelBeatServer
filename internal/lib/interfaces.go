package lib

type SpotifyApi interface {
	FetchPlaylistData(playlistId string, token string) (PlaylistData, error)
}
