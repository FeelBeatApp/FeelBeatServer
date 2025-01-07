package api

import (
	"net/http"
)

func SendJsonResponse(res *http.ResponseWriter, bytes []byte) error {
	(*res).Header().Set("Content-Type", "application/json")
	_, err := (*res).Write(bytes)

	return err
}
