package api

import (
	"encoding/json"
	"io"
)

func ParseBody(body io.ReadCloser, out any) error {
	bytes, err := ReadBody(body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, out)
}
