package spotify

import (
	"fmt"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
)

func GetUserId(token string) (string, error) {
	req, err := newGetApiCall("/me", token)
	if err != nil {
		return "", fmt.Errorf("Failed to create user id request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("User id request failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Auth request failed: %s", res.Status)
	}

	bytes, err := api.ReadBody(res.Body)
	if err != nil {
		return "", fmt.Errorf("Couldn't read user id request body: %w", err)
	}

	userid, err := jsonparser.GetString(bytes, "id")
	if err != nil {
		return "", fmt.Errorf("Failed to parse user profile: %w", err)
	}

	return userid, nil
}
