package feelbeaterror

import "net/http"

type ErrorCode string

const (
	Default = "unexpected_error"
)

func StatusCode(code ErrorCode) int {
	switch code {
	default:
		return http.StatusInternalServerError
	}
}
