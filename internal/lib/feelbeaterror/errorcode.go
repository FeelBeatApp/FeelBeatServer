package feelbeaterror

import "net/http"

type ErrorCode string

const (
	Default               = "Unexpected error occurred"
	AuthFailed            = "Authorization failed"
	LoadingPlaylistFailed = "Playlist loading failed"
)

func StatusCode(code ErrorCode) int {
	switch code {
	default:
		return http.StatusInternalServerError
	}
}
