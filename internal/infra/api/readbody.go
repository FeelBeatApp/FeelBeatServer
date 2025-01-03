package api

import (
	"io"
)

func ReadBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()

	bytes, err := io.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
