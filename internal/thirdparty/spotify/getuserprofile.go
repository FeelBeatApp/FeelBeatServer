package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
)

type profileResponse struct {
	Id     string `json:"id"`
	Name   string `json:"display_name"`
	Images []struct {
		Url string `json:"url"`
	} `json:"images"`
}

func GetUserProfile(token string) (lib.UserProfile, error) {
	req, err := newGetApiCall("/me", token)
	if err != nil {
		return lib.UserProfile{}, fmt.Errorf("Failed to create user id request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return lib.UserProfile{}, fmt.Errorf("User id request failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return lib.UserProfile{}, fmt.Errorf("Auth request failed: %s", res.Status)
	}

	bytes, err := api.ReadBody(res.Body)
	if err != nil {
		return lib.UserProfile{}, fmt.Errorf("Couldn't read user id request body: %w", err)
	}

	var response profileResponse
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return lib.UserProfile{}, fmt.Errorf("Failed to parse user profile: %w", err)
	}

	var imageUrl string
	if imageUrl = ""; len(response.Images) > 0 {
		imageUrl = response.Images[0].Url
	}

	return lib.UserProfile{
		Id:       response.Id,
		Name:     response.Name,
		ImageUrl: imageUrl,
	}, nil
}
