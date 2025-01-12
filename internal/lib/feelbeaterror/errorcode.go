package feelbeaterror

import "net/http"

type ErrorCode string

const (
	Default               = "Unexpected error occurred"
	AuthFailed            = "Authorization failed"
	LoadingPlaylistFailed = "Playlist loading failed"
	RoomNotFound          = "Room not found"
	RoomFull              = "Room is full"
	EncodingMessageFailed = "Encoding message failed"
)

func StatusCode(code ErrorCode) int {
	switch code {
	case RoomNotFound:
		return http.StatusNotFound
	case RoomFull:
		return http.StatusForbidden
	case AuthFailed:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
