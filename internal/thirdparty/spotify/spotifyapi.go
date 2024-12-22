package spotify

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/buger/jsonparser"
	"github.com/feelbeatapp/feelbeatserver/internal/lib"
)

func FetchSongDetails(spotifyId string, token string) {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", spotifyId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	title, err := jsonparser.GetString(bytes, "name")
	if err != nil {
		log.Fatal(err)

	}
	durationInMs, err := jsonparser.GetInt(bytes, "duration_ms")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(lib.SongDetails{
		Title:    title,
		Duration: time.Duration(durationInMs) * time.Millisecond,
	})
}
