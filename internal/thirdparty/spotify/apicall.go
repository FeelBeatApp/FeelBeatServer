package spotify

import "net/http"

func newGetApiCall(path string, token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", apiUrl+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}
