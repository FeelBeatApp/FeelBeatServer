package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/feelbeaterror"
)

type SpotifyApi struct {
}

func (s SpotifyApi) FetchPlaylistData(plalistId string, token string) (lib.PlaylistData, error) {
	url := fmt.Sprintf("/playlists/%s?additional_types=track&fields=name,images(url),tracks(items(track(id,images,name,artists(name),duration_ms)))", plalistId)
	req, err := newGetApiCall(url, token)
	if err != nil {
		return lib.PlaylistData{}, &feelbeaterror.FeelBeatError{
			DebugMessage: err.Error(),
			UserMessage:  feelbeaterror.LoadingPlaylistFailed,
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return lib.PlaylistData{}, &feelbeaterror.FeelBeatError{
			DebugMessage: err.Error(),
			UserMessage:  feelbeaterror.LoadingPlaylistFailed,
		}
	}

	bytes, err := api.ReadBody(res.Body)
	if err != nil {
		return lib.PlaylistData{}, &feelbeaterror.FeelBeatError{
			DebugMessage: err.Error(),
			UserMessage:  feelbeaterror.LoadingPlaylistFailed,
		}
	}
	var songsResponse playlistDataResponse
	err = json.Unmarshal(bytes, &songsResponse)
	if err != nil {
		return lib.PlaylistData{}, &feelbeaterror.FeelBeatError{
			DebugMessage: err.Error(),
			UserMessage:  feelbeaterror.LoadingPlaylistFailed,
		}
	}

	if len(songsResponse.Tracks.Items) == 0 {
		return lib.PlaylistData{}, &feelbeaterror.FeelBeatError{
			DebugMessage: "No songs in playlist",
			UserMessage:  feelbeaterror.LoadingPlaylistFailed,
		}
	}

	songs := make([]lib.Song, 0, len(songsResponse.Tracks.Items))
	for _, item := range songsResponse.Tracks.Items {
		artistNames := make([]string, 0, len(item.Track.Artists))
		for _, a := range item.Track.Artists {
			artistNames = append(artistNames, a.Name)
		}

		songs = append(songs, lib.Song{
			Id: item.Track.ID,
			Details: lib.SongDetails{
				Title:    item.Track.Name,
				Artist:   strings.Join(artistNames, " "),
				Duration: time.Duration(item.Track.DurationMs) * time.Millisecond,
			},
		})
	}

	var imageUrl string
	if imageUrl = ""; len(songsResponse.Images) > 0 {
		imageUrl = songsResponse.Images[0].Url
	}

	return lib.PlaylistData{
		Name:     songsResponse.Name,
		ImageUrl: imageUrl,
		Songs:    songs,
	}, nil
}
