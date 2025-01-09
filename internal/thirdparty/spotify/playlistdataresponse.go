package spotify

type playlistDataResponse struct {
	Name   string `json:"name"`
	Images []struct {
		Url string `json:"url"`
	} `json:"images"`
	Tracks struct {
		Items []struct {
			Track struct {
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Name       string `json:"name"`
				ID         string `json:"id"`
				DurationMs int    `json:"duration_ms"`
				Album      struct {
					Images []struct {
						Url string `json:"url"`
					} `json:"images"`
				} `json:"album"`
			} `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}
