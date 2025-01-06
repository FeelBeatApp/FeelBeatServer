package spotify

type playlistSongsResponse struct {
	Tracks struct {
		Items []struct {
			Track struct {
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Name       string `json:"name"`
				ID         string `json:"id"`
				DurationMs int    `json:"duration_ms"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}
